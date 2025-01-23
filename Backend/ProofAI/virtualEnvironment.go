package main

/*
 And in this file we have the virtualEnvironment function which is used to create a virtual environment and execute the model.
 And we have the downloadFile function which is used to download a file from the given URL. ( IPFS can be used to store the model and dataset )
 And we have the runCommand function which is used to run a command in the command prompt.
 And we have the ModelOuput struct which is used to parse the output of the Python model script.
 And we have the runPythonFile function which is used to run a Python file in the command prompt.
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
virtualEnvironment is a function to create a virtual environment and execute the model
 1. Create a virtual environment
 2. Activate the virtual environment
 3. Install the required packages
*/
func downloadFile(url string, filepath string, logger *log.Logger) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

/*
virtualEnvironment is a function to create a virtual environment and execute the model
 1. Create a virtual environment
 2. Activate the virtual environment
 3. Install the required packages
 4. Execute the model
 5. Return the model output
 6. logger is used to log the output of the commands in a file so that every transaction can be logged and give to the user
 7. You can train any model in the model.py file and use the trained model to predict the output.
    And you need to put all the requirements in the requirements.txt file. And the model.py file should be in the model directory.
    And for data you can use the dataset directory.
*/
func modelExecution() ([]byte, error) {

	timestamp := time.Now().Format("02_01_15_04_05") // day_month_hour_min_sec
	logFileName := fmt.Sprintf("%s_TransactionLog.txt", timestamp)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	logFile.Sync()
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	dirName := "virtualEnvironment"
	curr, _ := os.Getwd()
	dirPath := filepath.Join(ProofAI.modelExecutionDir, dirName)
	if _, err := os.Stat(dirPath); err == nil {
		if err := os.RemoveAll(dirPath); err != nil {
			logger.Printf("Failed to remove existing directory: %v", err)
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		logger.Printf("Error checking existing directory: %v", err)
		return nil, err
	}

	if err := os.Mkdir(dirPath, 0755); err != nil {
		logger.Printf("Failed to create directory: %v", err)
		return nil, err
	}

	if err := os.Chdir(dirPath); err != nil {
		logger.Printf("Failed to change directory: %v", err)
		return nil, err
	}
	logger.Printf("%s directory is created\n", dirName)

	runCommand("python -m venv Env", logger)
	runCommand(".\\Env\\Scripts\\activate", logger)
	runCommand(".\\Env\\Scripts\\activate && pip install -r ../model/requirements.txt", logger)

	model, err := runPythonFile(".\\Env\\Scripts\\activate && python ../model/model.py ../dataset/ ../model/knn_model.pkl", logger)

	if err != nil {
		logger.Printf("Failed to execute the Python model script: %v", err)
		return nil, err
	}

	// Return to the original working directory
	if err := os.Chdir(curr); err != nil {
		logger.Printf("Failed to change directory: %v", err)
		return nil, err
	}

	if err := os.RemoveAll(dirPath); err != nil {
		logger.Printf("Failed to remove directory: %v", err)
		return nil, err
	}

	return model, nil
}

/*
runCommand is a function to run a command in the command prompt
 1. command is the command to be executed
 2. logger is used to log the output of the command in a file so that every transaction can be logged and give to the user
 3. Return the output of the command
*/
func runCommand(command string, logger *log.Logger) string {
	cmd := exec.Command("cmd", "/C", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		logger.Printf("Command failed: %s\nError: %v\nOutput: %s", command, err, out.String())
	} else {
		logger.Printf("Command succeeded: %s\nOutput: %s", command, out.String())
	}
	return out.String()
}

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
runPythonFile is a function to run a Python file in the command prompt
 1. command is the command to be executed
 2. logger is used to log the output of the command in a file so that every transaction can be logged and give to the user
 3. Return the output of the command
*/
func runPythonFile(command string, logger *log.Logger) ([]byte, error) {
	cmd := exec.Command("cmd", "/C", command)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		logger.Printf("Command failed: %s\nError: %v\nOutput: %s", command, err, out.String())
		return nil, err
	}

	logger.Printf("Command succeeded: %s\n", command)

	var modelOutput ModelOuput
	err = json.Unmarshal(out.Bytes(), &modelOutput)
	if err != nil {
		logger.Printf("Failed to parse Python script output: %v\nOutput: %s\n", err, out.String())
		return nil, err
	}

	if modelOutput.Error != "" {
		logger.Printf("Python script error: %s\n", modelOutput.Error)
		return nil, fmt.Errorf("Python script error: %s", modelOutput.Error)
	}

	return out.Bytes(), nil
}
