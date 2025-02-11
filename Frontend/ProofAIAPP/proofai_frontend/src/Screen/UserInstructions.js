import React from "react";
import { ArrowLeft } from "lucide-react";

export default function UserInstructions() {
    const handleBack = () => {
        window.history.back();
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 p-2">
            <div className="mx-auto max-w-4xl rounded-lg bg-white/95 shadow-xl backdrop-blur">
                <div className="border-b border-slate-200 p-6">
                    <h1 className="text-2xl font-bold text-slate-800">
                        Model Training Instructions
                    </h1>
                </div>

                <div className="p-6">
                    <ol className="space-y-6 text-slate-700">
                        <li className="flex flex-col gap-2">
                            <p className="font-medium">Below is an example to create a directory structure:</p>
                            <div className="rounded-lg bg-slate-100 p-4 font-mono text-sm text-slate-800">
                                📂 project_root<br />
                                │── 📂 dataset<br />
                                │   ├── iris.data.txt  # Dataset file, you can use any dataset files and replace the name here<br />
                                │<br />
                                │── 📂 model<br />
                                │   ├── requirements.txt  # Dependencies, contain all the dependencies required for the model<br />
                                │   ├── model.py          # Model training script, you can place as much files as you want but make sure main file is model.py
                            </div>
                        </li>

                        <li className="flex flex-col gap-2">
                            <p className="font-medium">
                                Inside <span className="rounded bg-slate-200 px-1 py-0.5 font-mono">model/requirements.txt</span>, add the dependencies:
                            </p>
                            <div className="rounded-lg bg-slate-100 p-4 font-mono text-sm text-slate-800">
                                pandas<br />
                                requests<br />
                                scikit-learn<br />
                                joblib
                            </div>
                        </li>

                        <li className="flex flex-col gap-2">
                            <p className="font-medium">
                                In <span className="rounded bg-slate-200 px-1 py-0.5 font-mono">model/model.py</span>, include the following code:
                            </p>
                            <div className="max-h-96 overflow-y-auto rounded-lg bg-slate-100 p-4 font-mono text-sm text-slate-800">
                                <pre className="whitespace-pre-wrap">
                                    {`import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.neighbors import KNeighborsClassifier
import joblib  
import sys
import json
import os

# Get dataset path from command-line arguments
dataset_path = os.path.join(sys.argv[1], "iris.data.txt")
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
    joblib.dump(knn, model_filename)
    response = {"model_file": model_filename, "error": None}
    print(json.dumps(response))

except Exception as e:
    response = {"model_file": None, "error": str(e)}
    print(json.dumps(response))`}
                                </pre>
                            </div>
                        </li>
                    </ol>

                    <button
                        onClick={handleBack}
                        className="mt-6 flex items-center gap-2 rounded-lg bg-slate-800 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2"
                    >
                        <ArrowLeft className="h-4 w-4" />
                        Back
                    </button>
                </div>
            </div>
        </div>
    );
}