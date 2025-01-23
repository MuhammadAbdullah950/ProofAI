import React from "react";
import { FaDownload } from 'react-icons/fa';
import ReactJson from 'react-json-view';
import { useProofAiService } from "../ProofaiServiceContext";
import { useAlert } from "../Context/AlertContext";
import { useNavigate } from "react-router-dom";

const CurrentlyBlock = ({ }) => {

    const [showTransaction, setShowTransaction] = React.useState(false);
    const [block, setBlock] = React.useState([]);

    //  const alert = useAlert()
    const ProofAiService = useProofAiService()
    const handleBack = () => {
        window.history.back();
    };

    const getBlock = async () => {

        const response = await ProofAiService.getCurrentlyMinBlock()
        if (response.error) {
            return
        } else {
            if (response.block !== "null") {
                setBlock(response.block)
            } else {
                // make json object which shows no block is currently mining
                setBlock({ "Currently Mining Block": "No Block is currently mining" })
            }
        }
    }

    React.useEffect(() => {
        getBlock()
    }, [])




    return (
        <div style={styles.container}>
            <div style={styles.headerContainer}>
                <h1 style={styles.header}>Currently Mining Block</h1>
            </div>

            <div style={styles.buttonContainer}>
                <ReactJson
                    src={block}
                    theme="monokai"
                    collapsed={false}
                    displayDataTypes={false}
                />
            </div>
            <button style={{ ...styles.button, width: "10%", borderWidth: "8px", border: "2px solid white", color: "black", backgroundColor: "rgb(162, 168, 118)", paddingBottom: "20px" }} onClick={handleBack}>Back </button>
        </div>
    );
};

export default CurrentlyBlock;

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
