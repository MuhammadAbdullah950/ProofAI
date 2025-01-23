import React from "react";
import { useNavigate } from "react-router-dom";
import { useProofAiService } from "../ProofaiServiceContext";
import { useAlert } from "../Context/AlertContext";


// facing null error during getting blocks , need to send [] instead of null
const MinedTransactions = () => {
    const { showAlert } = useAlert();
    const navigate = useNavigate();
    const ProofAiService = useProofAiService();
    const [minedBlock, setMinedBlock] = React.useState([]);
    const [loading, setLoading] = React.useState(false);
    const [filterValue, setFilterValue] = React.useState("");

    const GetBlocks = async () => {
        setLoading(true);
        try {
            const response = await ProofAiService.getMinedBlocks(filterValue);

            if (response.error) {
                showAlert(response.error, "error");
            } else {
                if (response.blocks && Array.isArray(response.blocks)) {
                    // Ensure blocks is an array
                    setMinedBlock(response.blocks);
                } else {
                    // Set minedBlock to an empty array if blocks is null or not an array
                    setMinedBlock([]);
                    showAlert("No mined blocks found.", "success");
                }
            }
        } catch (error) {
            showAlert("Failed to fetch mined blocks.", "error");
            setMinedBlock([]); // Safely handle error by resetting minedBlock to an empty array
        } finally {
            setLoading(false);
        }
    };


    React.useEffect(() => {
        GetBlocks();
    }, []);

    const handleBack = () => {
        window.history.back();
    };

    const handleNavigate = (transactionString) => {
        const transaction = JSON.parse(transactionString); // Parse back to an object
        navigate(`/MinedTransaction`, { state: { transaction } }); // Pass transaction in state
    };

    const handleFilterBlocks = (e) => {
        setFilterValue(e.target.value);
        GetBlocks();
    };




    return (
        <div style={styles.container}>
            <div style={styles.headerContainer}>
                <h1 style={styles.header}>Mined Blocks</h1>
                <select style={styles.dropdown} onChange={GetBlocks} >
                    <option value="">All Transactions</option>
                    <option value="own">Own Transactions</option>
                </select>
            </div>

            {minedBlock && minedBlock.length > 0 ? (
                minedBlock.map((block) => (
                    <select
                        key={block.BlockNum}
                        style={styles.Button}
                        onChange={(e) => handleNavigate(e.target.value)}
                    >
                        <option value="">
                            BlockNum: {block.blockNum}
                        </option>
                        {block.transactions && block.transactions.length > 0 ? (
                            block.transactions.map((transaction) => (
                                <option
                                    key={transaction.transactionId}
                                    value={JSON.stringify(transaction)}
                                >
                                    Transaction = Nonce: {transaction.nonce}, From: {transaction.from}
                                </option>
                            ))
                        ) : (
                            <option disabled>No transactions</option>
                        )}
                    </select>
                ))
            ) : (
                <p style={{ textAlign: "center", color: "gray" }}>
                    No mined blocks found.
                </p>
            )}

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
};

export default MinedTransactions;

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
    },
    dropdown: {
        padding: "8px 12px",
        border: "2px solid rgb(190, 206, 70)",
        borderRadius: "5px",
        backgroundColor: "#fff",
        color: "#333",
        fontSize: "14px",
        cursor: "pointer",
    },
    buttonContainer: {
        display: "flex",
        flexDirection: "column",
        gap: "10px",
        overflowY: "auto",
        maxHeight: "470px",
        marginBottom: "20px",
        backgroundColor: "rgb(124, 124, 124)",
        padding: "10px",
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
    Button: {
        padding: "8px 16px",
        border: "2px solid #595f66",
        borderRadius: "5px",
        backgroundColor: "#333",
        color: "#fff",
        cursor: "pointer",
        fontSize: "14px",
        transition: "background-color 0.3s ease",
    },
};
