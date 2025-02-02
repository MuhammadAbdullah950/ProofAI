package main

/*		In this file we run the model in a virtual environment.
1-		ModelOuput is a struct to parse the output of the Python model script.
2-		downloadFromIPFS function which is used to download files from the given URL. ( IPFS can be used to store the model and dataset )
3-		modelExecution function which is used to create a virtual environment and execute the model.
4-    	readLogFileToBytes function which is used to read the log file into a byte array.
5-		changeDir function which is used to change the current working directory.
6-		runCommand function which is used to run a command in the command prompt.
7-		runPythonFile function which is used to run a Python file in the command prompt.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

/*
ModelOuput is a struct to parse the output of the Python model script
 1. Model is the output of the model
 2. Error is the error message from the model
*/
type ModelOuput struct {
	Model string `json:"model"`
	Error string `json:"error"`
}

/*
downloadFromIPFS is a function to download files from IPFS
 1. cid: content identifier of the file
 2. outputDir: output directory to save the files
    Create the output directory
    Prepare the request
    Check the status code
    Parse the response
    Save each file
*/
func downloadFromIPFS(cid string, outputDir string) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	serviceMachineURl := "http://" + ProofAI.selfMiningDetail.serviceMachineAddr

	data := bytes.NewBufferString(fmt.Sprintf("cid=%s", cid))
	resp, err := http.Post(serviceMachineURl+"/fetch", "application/x-www-form-urlencoded", data)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	// Parse response
	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Message)
	}

	for _, file := range response.Files {
		filePath := filepath.Join(outputDir, file.Name)

		// Write the content directly without base64 decoding
		err = os.WriteFile(filePath, []byte(file.Content), 0644)
		if err != nil {
			return fmt.Errorf("failed to save file %s: %v", file.Name, err)
		}
	}

	return nil
}

/*
modelExecution is a function to create a virtual environment and execute the model
*/
func modelExecution(CID_Input_dataSet string, CID_Input_model string, dirPath string) ([]byte, []byte, error) {

	// create log file
	timestamp := time.Now().Format("02_01_15_04_05")
	logFileName := fmt.Sprintf("%s_TransactionLog.txt", timestamp)
	currentdir, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	logFilePath := filepath.Join(currentdir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create or open log file: %w", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", 0)

	// Download the dataset and model from IPFS
	err = downloadFromIPFS(CID_Input_dataSet, filepath.Join(dirPath, "dataset"))
	if err != nil {
		logger.Printf("Error downloading dataset from IPFS: %v", err)
		logData, err := readLogFileToBytes(logFilePath)
		if err != nil {
			logger.Printf("Error reading log file to bytes: %v", err)
			return nil, nil, err
		}
		return nil, logData, err
	}
	err = downloadFromIPFS(CID_Input_model, filepath.Join(dirPath, "model"))
	if err != nil {
		logger.Printf("Error downloading model from IPFS: %v", err)
		logData, err := readLogFileToBytes(logFilePath)
		if err != nil {
			logger.Printf("Error reading log file to bytes: %v", err)
			return nil, nil, err
		}
		return nil, logData, err
	}

	virtualEnvDir := filepath.Join(dirPath, "virtualEnvironment")

	if err := os.Mkdir(virtualEnvDir, 0777); err != nil {
		logger.Printf("Failed to create directory: %v", err)
		return nil, nil, err
	}
	if err := os.Chdir(virtualEnvDir); err != nil {
		logger.Printf("Failed to change directory: %v", err)
		return nil, nil, err
	}
	defer func() {
		if err := os.Chdir(currentdir); err != nil {
			logger.Printf("Failed to return to original directory: %v", err)
		}
	}()
	//	logger.Printf("Directory %s is created and in use\n", virtualEnvDir)

	// Execute commands for virtual environment setup and model execution
	if _, err := runCommand("python -m venv Env"); err != nil {
		logger.Printf("Failed to create virtual environment : %v", err)
		logData, err := readLogFileToBytes(logFilePath)
		if err != nil {
			logger.Printf("Error reading log file to bytes: %v", err)
			return nil, nil, err
		}
		return nil, logData, err
	}
	logger.Printf("Virtual environment created\n")

	if _, err := runCommand(".\\Env\\Scripts\\activate"); err != nil {
		logger.Printf("Failed to activate virtual environment: %v", err)
		logData, err := readLogFileToBytes(logFilePath)
		if err != nil {
			logger.Printf("Error reading log file to bytes: %v", err)
			return nil, nil, err
		}
		return nil, logData, err
	}
	logger.Printf("Virtual environment activated\n")

	if _, err := runCommand(".\\Env\\Scripts\\activate && pip install -r ../model/requirements.txt"); err != nil {
		logger.Printf("Failed to install required packages: %v", err)
		logData, err := readLogFileToBytes(logFilePath)
		if err != nil {
			logger.Printf("Error reading log file to bytes: %v", err)
			return nil, nil, err
		}
		return nil, logData, err
	}
	logger.Printf("Required packages installed\n")

	// Execute the Python model script
	model, err := runPythonFile(".\\Env\\Scripts\\activate && python ../model/model.py ../dataset/ ../model/knn_model.pkl")
	if err != nil {
		logger.Printf("Failed to execute the Python model script: %v", err)
		logData, err := readLogFileToBytes(logFilePath)
		if err != nil {
			logger.Printf("Error reading log file to bytes: %v", err)
			return nil, nil, err
		}
		return nil, logData, err
	}

	logger.Printf("Model executed successfully\n")

	// Log file debugging

	// Read the log file into a byte array
	logData, err := readLogFileToBytes(logFilePath)
	if err != nil {
		logger.Printf("Error reading log file to bytes: %v", err)
		return nil, nil, err
	}

	return model, logData, nil
}

/* Function to read log file into a byte array
 */
func readLogFileToBytes(logFilePath string) ([]byte, error) {
	// Open the log file in read-only mode
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Read the file content into a byte array
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	return content, nil
}

/*
changeDir is a function to change the current working directory
*/
func changeDir(dirPath string) error {
	if err := os.Chdir(dirPath); err != nil {
		return err
	}
	return nil
}

/*
runCommand is a function to run a command in the command prompt
 1. command is the command to be executed
 2. logger is used to log the output of the command in a file so that every transaction can be logged and give to the user
 3. Return the output of the command
*/
func runCommand(command string) (string, error) {

	cmd := exec.Command("cmd", "/C", command)

	var out bytes.Buffer
	var stderr bytes.Buffer

	// Capture both stdout and stderr
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("failed to execute command: %s, error: %w, stderr: %s", command, err, stderr.String())
	}

	return out.String(), nil
}

/*
runPythonFile is a function to run a Python file in the command prompt
 1. command is the command to be executed
 2. logger is used to log the output of the command in a file so that every transaction can be logged and give to the user
 3. Return the output of the command
*/
func runPythonFile(command string) ([]byte, error) {
	cmd := exec.Command("cmd", "/C", command)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Command failed: %s\nError: %v\nOutput: %s", command, err, out.String())
	}

	var modelOutput ModelOuput
	err = json.Unmarshal(out.Bytes(), &modelOutput)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Python script output: %v\nOutput: %s\n", err, out.String())
	}

	if modelOutput.Error != "" {
		return nil, fmt.Errorf("Python script error: %s", modelOutput.Error)
	}

	return out.Bytes(), nil
}
