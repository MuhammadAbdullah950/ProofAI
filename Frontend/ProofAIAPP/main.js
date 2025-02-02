const { spawn } = require('child_process');
const { app, ipcMain, BrowserWindow, dialog } = require('electron');
const path = require('path');
const fs = require('fs');

let mainWindow;
let backendProcess;
let isQuitting = false;

ipcMain.on("restart-app", () => {
  app.relaunch();
  app.quit();
});

function createWindow() {
  const win = new BrowserWindow({

    minHeight: 600,
    minWidth: 900,

    width: 1200,
    height: 750,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
      webSecurity: false,
      devTools: true
    }
  });

  win.setMenuBarVisibility(false);
  win.setTitle('ProofAI');

  if (!app.isPackaged) {
    win.loadFile('Electron-Framework/waiting.html');

    const pollServer = () => {
      require('http').get('http://localhost:3000', (res) => {
        if (res.statusCode === 200) {
          win.loadURL('http://localhost:3000');
        } else {
          setTimeout(pollServer, 1000);
        }
      }).on('error', () => {
        setTimeout(pollServer, 1000);
      });
    };
    pollServer();
  } else {
    const htmlPath = path.join(app.getAppPath(), 'proofai_frontend', 'build', 'index.html');
    console.log('Loading frontend from:', htmlPath);
    win.loadFile(htmlPath);
  }

  return win;
}

function startBackend() {

  let backendPath;
  let backendOutput = '';

  if (app.isPackaged) {
    const possiblePaths = [
      path.join(process.resourcesPath, 'ProoAiBackend.exe'),
      path.join(app.getAppPath(), 'ProoAiBackend.exe'),
      path.join(process.resourcesPath, 'app', 'ProoAiBackend.exe'),
      path.join(app.getPath('exe'), '..', 'ProoAiBackend.exe')
    ];

    for (const testPath of possiblePaths) {
      console.log('Checking backend path:', testPath);
      if (fs.existsSync(testPath)) {
        backendPath = testPath;
        break;
      }
    }
  } else {
    backendPath = path.join(__dirname, 'ProoAiBackend.exe');
  }

  if (!backendPath || !fs.existsSync(backendPath)) {
    const error = 'Backend executable not found';
    console.error(error);
    dialog.showErrorBox('Backend Error', error);
    return null;
  }

  const logPath = path.join(app.getPath('userData'), 'backend.log');
  const logStream = fs.createWriteStream(logPath, { flags: 'a' });

  const proc = spawn(backendPath, [], {
    stdio: 'pipe',
    detached: false,
    env: {
      ...process.env,
      ELECTRON_BACKEND: '1',
      PATH: process.env.PATH
    },
    cwd: path.dirname(backendPath)
  });

  logStream.write(`Started backend process with PID: ${proc.pid}\n`);

  proc.stdout.on('data', (data) => {
    const output = data.toString();
    backendOutput += output;
    console.log('Backend stdout:', output);
    logStream.write(`stdout: ${output}\n`);
  });

  proc.stderr.on('data', (data) => {
    const output = data.toString();
    backendOutput += output;
    console.error('Backend stderr:', output);
    logStream.write(`stderr: ${output}\n`);
  });

  proc.on('error', (err) => {
    console.error('Backend process error:', err);
    logStream.write(`Process error: ${err.message}\n`);
    if (!isQuitting) {
      dialog.showErrorBox('Backend Error', `Backend error: ${err.message}\n\nOutput:\n${backendOutput}`);
    }
  });

  proc.on('exit', (code, signal) => {
    console.log(`Backend process exited with code ${code} and signal ${signal}`);
    logStream.write(`Process exited with code ${code} and signal ${signal}\n`);

    if (!isQuitting && code !== 0) {
      dialog.showErrorBox(
        'Backend Error',
        `Backend process exited unexpectedly.\nExit code: ${code}\nSignal: ${signal}\n\nOutput:\n${backendOutput}`
      );

      setTimeout(() => {
        if (!isQuitting) {
          console.log('Attempting to restart backend...');
          startBackend();
        }
      }, 1000);
    }
  });

  return proc;
}

app.whenReady().then(() => {
  console.log('App directory:', app.getAppPath());
  console.log('Resource path:', process.resourcesPath);

  try {
    backendProcess = startBackend();
    if (backendProcess) {
      mainWindow = createWindow();
    } else {
      dialog.showErrorBox('Error', 'Failed to start backend process');
      app.quit();
    }
  } catch (error) {
    console.error('Failed to initialize application:', error);
    dialog.showErrorBox('Error', `Failed to initialize application: ${error.message}`);
    app.quit();
  }
});

app.on('window-all-closed', () => {
  isQuitting = true;
  if (backendProcess) {
    try {
      if (process.platform === 'win32') {
        spawn('taskkill', ['/pid', backendProcess.pid, '/f', '/t']);
      } else {
        backendProcess.kill();
      }
    } catch (error) {
      console.error('Error killing backend process:', error);
    }
  }

  if (process.platform !== 'darwin') {
    app.quit();
  }
});

process.on('uncaughtException', (error) => {
  console.error('Uncaught exception:', error);
  if (!isQuitting) {
    dialog.showErrorBox('Error', `Uncaught exception: ${error.message}`);
  }
});

process.on('unhandledRejection', (error) => {
  console.error('Unhandled rejection:', error);
  if (!isQuitting) {
    dialog.showErrorBox('Error', `Unhandled rejection: ${error.message}`);
  }
});