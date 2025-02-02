import React from "react";
import { FaDownload } from 'react-icons/fa';
import ReactJson from 'react-json-view';
import { useLocation } from "react-router-dom";

const MinedTransaction = ({ }) => {

    const [showTransaction, setShowTransaction] = React.useState(false);
    const location = useLocation();
    const { transaction } = location.state;
    const [isDiabled, setisDiabled] = React.useState(false);

    React.useEffect(() => {
        if (transaction?.type === "block") {
            setisDiabled(true);
        }
    }, [transaction])

    const handleBack = () => {
        window.history.back();
    };

    const handleModelDownload = () => {
        try {
            if (!transaction.model_output) {
                alert("Model output is missing.");
                return;
            }
            const jsonString = atob(transaction.model_output);
            let jsonData;
            try {
                jsonData = JSON.parse(jsonString);
            } catch (e) {
                alert("Error parsing JSON: " + e.message);
                return;
            }

            if (jsonData.error) {
                alert("Error in model file: " + jsonData.error);
                return;
            }

            if (!jsonData.model) {
                alert("Model data is missing.");
                return;
            }
            let binaryString;
            try {
                binaryString = atob(jsonData.model);
            } catch (e) {
                alert("Error decoding model data: " + e.message);
                return;
            }

            const byteArray = Uint8Array.from(binaryString, c => c.charCodeAt(0));
            const blob = new Blob([byteArray], { type: 'application/octet-stream' });
            const url = window.URL.createObjectURL(blob);

            const a = document.createElement('a');
            a.href = url;
            a.download = 'modeFile.pkl';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } catch (e) {
            alert("An unexpected error occurred: " + e.message);
            console.error("Unexpected error:", e);
        }
    };


    const handlTransactionLogDownload = () => {
        // Download the transaction data
        const modelData = atob(transaction.transactionLog); // Decode base64 string
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
        a.download = 'TransactionLogs.txt'; // Specify the file name here
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
                    <label style={{ ...styles.dropdown, pointerEvents: isDiabled ? "none" : "auto", opacity: isDiabled ? 0.5 : 1 }} onClick={!isDiabled ? handleModelDownload : undefined}>
                        Download Trained Model <FaDownload />
                    </label>
                    <label style={{ ...styles.dropdown, pointerEvents: isDiabled ? "none" : "auto", opacity: isDiabled ? 0.5 : 1 }} onClick={!isDiabled ? handlTransactionLogDownload : undefined}>
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
        </div >
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
