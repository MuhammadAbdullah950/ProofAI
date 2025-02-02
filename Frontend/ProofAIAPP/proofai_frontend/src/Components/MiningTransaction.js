import React from "react";
import ReactJson from 'react-json-view';
import { ClipLoader } from 'react-spinners';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from '../Context/AlertContext';
import { FaCheckCircle } from 'react-icons/fa';
import { useEffect } from "react";

const MiningTransaction = ({ setmineTransaction, transactionJson }) => {
    const handleBack = () => {
        setmineTransaction(false);
    };

    const [isTransactionMined, setisTransactionMined] = React.useState(false);
    const [intervalId, setIntervalId] = React.useState(null);
    const proofAiService = useProofAiService();

    const chkTransactionStatus = async () => {
        const response = await proofAiService.transactionConfirmation(transactionJson.from, transactionJson.nonce);
        if (response.transaction === "Confirmed") {
            setisTransactionMined(true);
        } else {
            alert("Transaction is not mined yet")
        }
    }

    return (
        <div style={styles.container}>
            <h1 style={styles.label}>Mining Transaction</h1>
            {!isTransactionMined &&
                <>
                    <div style={styles.loaderContainer}>
                        <ClipLoader color={'#123abc'} loading={true} size={50} />
                    </div>
                    <div style={styles.jsonBox}>
                        <ReactJson
                            src={transactionJson}
                            theme="monokai"
                            collapsed={false}
                            displayDataTypes={false}
                            displayObjectSize={true}
                        />
                    </div>
                </>
            }
            {isTransactionMined && <FaCheckCircle color="green" size={20} />}
            <div style={{ display: "flex", gap: "10px" }} >
                <button style={styles.button} onClick={chkTransactionStatus}>Check Status</button>
                <button style={styles.button} onClick={handleBack}>Back</button>
            </div>
        </div>
    );
};

export default MiningTransaction;

const styles = {
    container: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '10px',
        padding: '10px',
        backgroundColor: "rgb(152, 156, 158)",
        borderRadius: '10px',
        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
        width: '80%',
        maxWidth: '1200px',
        margin: 'auto',
        marginTop: '60px',
    },
    label: {
        fontSize: '27px',
        fontWeight: 'bold',
        color: '#333',
        textAlign: 'center',
    },
    loaderContainer: {
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '50px',
    },
    jsonBox: {
        width: '100%',
        height: '300px',
        padding: '10px',
        backgroundColor: '#282c34',
        borderRadius: '5px',
        boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
        overflowY: 'auto',
        overflowX: 'auto',
        wordWrap: 'break-word',
    },
    button: {
        padding: '10px 20px',
        border: 'none',
        borderRadius: '5px',
        backgroundColor: 'rgb(15, 67, 89)',
        color: '#fff',
        cursor: 'pointer',
        fontSize: '14px',
        fontWeight: 'bold',
        transition: 'background-color 0.3s ease',
        border: '2px solid rgb(54, 53, 50)',
    },
};
