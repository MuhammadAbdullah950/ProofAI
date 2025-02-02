import React from "react";

export default function UserInstructions() {

    const handleBack = () => {
        window.history.back();
    }
    return (
        <div className="p-6 bg-gray-100 rounded-lg shadow-md" style={styles.container} >
            <h1 className="text-xl font-bold mb-4"> Model Training Instructions</h1>
            <ol className="list-decimal list-inside space-y-2">
                <li>Below is an example to create a directory structure.</li>
                <pre className="bg-gray-200 p-4 rounded-md text-sm">
                    {`📂 project_root
│── 📂 dataset  
│   ├── iris.data.txt  # Dataset file  , you can use any dataset files and replace the name here 
│
│── 📂 model
│   ├── requirements.txt  # Dependencies , contain all the dependencies required for the model
│   ├── model.py          # Model training script, you can place as much files as you want but make sure main file is model.py`}
                </pre>
                <li>
                    Inside <strong>model/requirements.txt</strong>, add the dependencies:
                    <pre className="bg-gray-200 p-4 rounded-md text-sm">{`
pandas
requests
scikit-learn
joblib`}</pre>
                </li>
                <li>
                    In <strong>model/model.py</strong>, include the following code:
                    <pre className="bg-gray-200 p-4 rounded-md text-sm overflow-x-auto">{`
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.neighbors import KNeighborsClassifier
import joblib  
import sys
import json
import os

# Get dataset path from command-line arguments
dataset_path = os.path.join(sys.argv[1], "iris.data.txt") # sys.argv[1] is the dataset directory path and you can replace the file name with your dataset file name
model_filename = "modelFile.pkl"

try:
    # Load the dataset
    column_names = ["SepalLength", "SepalWidth", "PetalLength", "PetalWidth", "Species"]
    df = pd.read_csv(dataset_path, header=None, names=column_names)

    # Convert species names to numbers
    df["Species"] = df["Species"].map({
        "Iris-setosa": 0,
        "Iris-versicolor": 1,
        "Iris-virginica": 2
    })

    # Ensure there are no missing values
    if df.isnull().sum().any():
        raise ValueError("Dataset contains missing values. Please clean the data and try again.")

    # Split data into features and labels
    X = df.drop(columns=["Species"])
    y = df["Species"]

    # Split into training and testing sets
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    # Train KNN Model
    knn = KNeighborsClassifier(n_neighbors=3)
    knn.fit(X_train, y_train)

    # Save trained model
    joblib.dump(knn, model_filename)                            # necessary Step
    response = {"model_file": model_filename, "error": None}    # necessary Step
    print(json.dumps(response))                                 # necessary Step

except Exception as e:
    response = {"model_file": None, "error": str(e)}  # necessary Step
    print(json.dumps(response))                   # necessary Step
          `}</pre>
                </li>
            </ol>

            <button
                style={{
                    ...styles.button,
                    width: "10%",
                    borderWidth: "8px",
                    border: "2px solid white",
                    color: "black",
                    backgroundColor: "rgb(162, 168, 118)",
                }}
                onClick={handleBack}
            >
                Back
            </button>
        </div>
    );
}

const styles = {

    container: {
        display: "flex",
        flexDirection: "column",
        width: "100%",
        minHeight: "100vh",
        backgroundColor: "#f4f6f7",
        boxSizing: "border-box",
        padding: "20px",
    },

    button: {
        padding: "12px",
        border: "2px solid rgb(57, 96, 111)",
        borderRadius: "5px",
        color: "white",
        backgroundColor: "rgb(41, 41, 43)",
        cursor: "pointer",
        fontSize: "16px",
        fontWeight: "bold",
        transition: "all 0.3s ease",
    },
}
