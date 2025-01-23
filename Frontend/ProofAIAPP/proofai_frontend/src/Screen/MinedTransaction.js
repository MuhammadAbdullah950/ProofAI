import React from "react";
import { FaDownload } from 'react-icons/fa';
import ReactJson from 'react-json-view';
import { useLocation } from "react-router-dom";

const MinedTransaction = ({ }) => {

    const [showTransaction, setShowTransaction] = React.useState(false);
    const location = useLocation();
    const { transaction } = location.state; // Access the passed transaction


    const handleBack = () => {
        window.history.back();
    };

    const handleDownload = () => {
        // Download the transaction data
        const modelData = atob(transaction.model_output); // Decode base64 string
        const byteArray = new Uint8Array(modelData.length);

        // Convert string to byte array
        for (let i = 0; i < modelData.length; i++) {
            byteArray[i] = modelData.charCodeAt(i);
        }

        // Create a Blob from the byte array and download it
        const blob = new Blob([byteArray], { type: 'application/octet-stream' });
        const url = window.URL.createObjectURL(blob);

        // Create a temporary anchor element to trigger the download
        const a = document.createElement('a');
        a.href = url;
        a.download = 'model_file.pkl'; // Specify the file name here
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    };



    return (
        <div style={styles.container}>
            <div style={styles.headerContainer}>
                <h1 style={styles.header}>Mined Transactions</h1>
                <div style={{ display: "flex", gap: "10px" }}>
                    <label style={styles.dropdown} onClick={handleDownload}>
                        Download Trained Model <FaDownload />
                    </label>
                    <label style={styles.dropdown} onClick={handleDownload}>
                        Download Transaction Logs <FaDownload />
                    </label>

                </div>

            </div>

            <div style={styles.buttonContainer}>
                <ReactJson
                    src={transaction}
                    theme="monokai"
                    collapsed={false}
                    displayDataTypes={false}
                />
            </div>
            <button style={{ ...styles.button, width: "10%", borderWidth: "8px", border: "2px solid white", color: "black", backgroundColor: "rgb(162, 168, 118)", paddingBottom: "20px" }} onClick={handleBack}>Back </button>
        </div>
    );
};

export default MinedTransaction;

const styles = {
    container: {
        display: "flex",
        flexDirection: "column",
        width: "100%",
        minHeight: "100vh",
        backgroundColor: "#f4f6f7",
        boxSizing: "border-box",
        padding: "20px 20px",

    },
    headerContainer: {
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        marginBottom: "20px",
    },
    header: {
        fontSize: "28px",
        fontWeight: "bold",
        color: "#333",
        textAlign: "left",
    },
    dropdown: {
        padding: "8px 12px",
        border: "2px solid rgb(190, 206, 70)",
        borderRadius: "5px",
        backgroundColor: "#fff",
        color: "#333",
        fontSize: "14px",
        cursor: "pointer",
        outline: "none",
    },
    buttonContainer: {
        display: "flex",
        flexDirection: "column",
        gap: "10px",
        overflowY: "auto",
        maxHeight: "470px",
        marginBottom: "5px",
        backgroundColor: "rgb(124, 124, 124)",
        padding: "10px 10px",

    },
    button: {
        padding: "12px",
        border: "2px solid rgb(57, 96, 111)",
        width: "100vh",
        borderRadius: "2px",
        color: "white",
        backgroundColor: "rgb(41, 41, 43)",
        cursor: "pointer",
        fontSize: "16px",
        fontWeight: "bold",
        transition: "all 0.3s ease",
    },
    backButton: {
        padding: "10px 20px",
        marginTop: "20px",
        border: "2px solid #595f66",
        borderRadius: "8px",
        backgroundColor: "#555",
        color: "#fff",
        fontSize: "14px",
        fontWeight: "bold",
        cursor: "pointer",
        transition: "all 0.3s ease",
    },
};
