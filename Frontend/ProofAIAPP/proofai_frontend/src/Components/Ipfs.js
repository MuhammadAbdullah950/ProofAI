import React, { useState } from 'react';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from "../Context/AlertContext";
import { ClipLoader } from 'react-spinners';
import axios from 'axios';

const IPFS = ({ setipfsCont }) => {
    const [uploadResponse, setUploadResponse] = useState("");
    const [ipfsContent, setIpfsContent] = useState(false);
    const [waitingIcon, setWaitingIcon] = useState(false);

    const ProofAiService = useProofAiService();
    const { showAlert } = useAlert();
    const handleBackButton = () => setipfsCont(false);

    const handleUpload = async (event) => {
        setIpfsContent(false);
        setUploadResponse("");
        const files = event.target.files;
        if (!files.length) return;
        const data = new FormData();
        Array.from(files).forEach(file => data.append("files", file));
        setWaitingIcon(true);

        try {
            const response = await ProofAiService.uploadDataOnIPfs(data);
            setWaitingIcon(false);

            if (response.error) {
                showAlert(`Error during upload: ${response.error}`, "error");
            } else {
                setUploadResponse(response);
            }
        } catch (error) {
            setWaitingIcon(false);
            alert(error)
            showAlert("An unexpected error occurred during upload.", "error");
        }
    };

    const handleCopyCID = () => {
        navigator.clipboard.writeText(uploadResponse);
        showAlert("CID copied to clipboard!", "success");
    };

    return (
        <div style={styles.box}>
            <label style={styles.label}>Upload Folder to IPFS</label>
            <input
                type="file"
                id="file"
                style={styles.input}
                webkitdirectory="true"
                directory="true"
                multiple
                onChange={handleUpload}
            />
            <button style={styles.button}>Upload</button>

            {waitingIcon && <ClipLoader color="#123abc" loading={true} size={50} />}

            {uploadResponse && (
                <div style={styles.cidBox}>
                    <label>CID of Uploaded Folder:</label>
                    <div style={styles.cidContainer}>
                        <span style={styles.cidText}>{uploadResponse}</span>
                        <button style={styles.copyButton} onClick={handleCopyCID}>
                            Copy
                        </button>
                    </div>
                </div>
            )}

            <button style={{ ...styles.button, width: "30%" }} onClick={handleBackButton}>Back</button>
        </div>
    );
};

export default IPFS;

const styles = {
    label: {
        fontSize: "1.5rem",
        fontWeight: "bold",
        color: "#333",
        textAlign: "center",
        marginBottom: "10px",
    },
    box: {
        display: "flex",
        flexDirection: "column",
        gap: "15px",
        padding: "20px",
        backgroundColor: "#989c9e",
        borderRadius: "10px",
        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
        width: "50%",
        margin: "0 auto",
    },
    input: {
        display: "block",
        width: "100%",
        padding: "10px",
        borderRadius: "5px",
        border: "1px solid #ccc",
        boxSizing: "border-box",
    },
    button: {
        padding: "10px 20px",
        border: "none",
        borderRadius: "5px",
        backgroundColor: "#454f53",
        color: "#fff",
        cursor: "pointer",
        fontSize: "14px",
        fontWeight: "bold",
        transition: "background-color 0.3s ease",
        border: "2px solid #080807",
    },
    cidBox: {
        padding: "10px",
        backgroundColor: "#fff",
        borderRadius: "5px",
        border: "1px solid #ccc",
        boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)",
    },
    cidContainer: {
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        marginTop: "5px",
    },
    cidText: {
        fontSize: "14px",
        color: "#333",
        wordBreak: "break-word",
    },
    copyButton: {
        padding: "5px 10px",
        fontSize: "12px",
        color: "#fff",
        backgroundColor: "#4CAF50",
        border: "none",
        borderRadius: "5px",
        cursor: "pointer",
        fontWeight: "bold",
        transition: "background-color 0.3s ease",
    },
};
