import React from "react";
import { FaDownload } from 'react-icons/fa';
import ReactJson from 'react-json-view';

const useScreenMinedTransaction = ({ transaction, showTransaction, setShowTransaction, handleBack }) => {



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
        const modelData = atob(transaction.transactionLog); // Decode base64 string
        const byteArray = new Uint8Array(modelData.length);

        for (let i = 0; i < modelData.length; i++) {
            byteArray[i] = modelData.charCodeAt(i);
        }

        const blob = new Blob([byteArray], { type: 'application/octet-stream' });
        const url = window.URL.createObjectURL(blob);

        const a = document.createElement('a');
        a.href = url;
        a.download = 'TransactionLogs.txt';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    };

    return {
        showTransaction,
        setShowTransaction,
        handleBack,
        handleModelDownload,
        handlTransactionLogDownload,
    }
};

export default useScreenMinedTransaction;
