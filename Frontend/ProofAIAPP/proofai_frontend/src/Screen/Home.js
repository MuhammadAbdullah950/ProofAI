import React, { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useProofAiService } from "../ProofaiServiceContext";


const Home = () => {
    const navigate = useNavigate()

    const handleNewTransaction = () => {
        navigate('/NewTransaction')
    }

    const handleMinedTransactions = () => {
        navigate('/MinedTransactions')
    }

    const ProofAiService = useProofAiService()

    const handleMiningTransaction = async () => {
        navigate('/CurrentlyBlock')
    }

    return (

        <div style={styles.container} >
            <div style={styles.maxWidth}>
                <div style={{ ...styles.boxShadow, ...styles.bgGray }}>
                    <div style={styles.padding}>
                        <h2 style={{ ...styles.textCenter, ...styles.textWhite, ...styles.text3xl, ...styles.fontExtrabold }}>Transaction Dashboard</h2>

                        <div style={{ ...styles.roundedShadow, ...styles.mt4, display: "flex", flexDirection: "column", gap: "10px" }}>
                            <button style={styles.input} onClick={handleNewTransaction}  >New Transaction</button>
                            <button style={styles.input} onClick={handleMinedTransactions} >Mined Blocks</button>
                            <button style={styles.input} onClick={handleMiningTransaction} >Currently Mining block</button>
                        </div>
                    </div>
                </div>
            </div >
        </div>
    );
}

export default Home;


const styles = {

    container: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: "100%",
        height: "100vh",
        backgroundColor: "#f4f6f7",
        flexDirection: 'column',
    },

    label: {
        fontSize: '0.975rem',
        lineHeight: '1.25rem',
        color: '#cbd5e0',
        marginbottom: '0.5rem',
    },

    maxWidth: {
        maxWidth: '32rem',
        width: '100%',
    },
    boxShadow: {
        boxShadow: '0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
    },
    bgGray: {
        backgroundColor: '#2d3748',
        borderRadius: '0.5rem',
        overflow: 'hidden',
    },
    padding: {
        padding: '2rem',
    },
    textCenter: {
        textAlign: 'center',
    },
    textWhite: {
        color: '#fff',
    },
    text3xl: {
        fontSize: '1.875rem',
        lineHeight: '2.25rem',
    },
    fontExtrabold: {
        fontWeight: '800',
    },
    mt4: {
        marginTop: '1rem',
    },
    textGray: {
        color: '#a0aec0',
    },
    mt8: {
        marginTop: '2rem',
    },
    roundedShadow: {
        borderRadius: '0.375rem',
        boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
    },
    srOnly: {
        position: 'absolute',
        width: '1px',
        height: '1px',
        padding: '0',
        margin: '-1px',
        overflow: 'hidden',
        clip: 'rect(0, 0, 0, 0)',
        border: '0',
    },
    input: {
        appearance: 'none',
        display: 'block',
        width: '100%',
        padding: '0.75rem',
        border: '1px solid #4a5568',
        backgroundColor: '#4a5568',
        color: '#fff',
        borderRadius: '0.375rem',
        fontSize: '15px',

        // space bwteen character
        letterSpacing: '0.5px',
        border: '1px solid #bc8c49',

    },
    flex: {
        display: 'flex',
    },
    itemsCenter: {
        alignItems: 'center',
    },
    justifyBetween: {
        justifyContent: 'space-between',
    },
    checkbox: {
        height: '1rem',
        width: '1rem',
        color: '#667eea',
        borderColor: '#4a5568',
        borderRadius: '0.25rem',
    },
    ml2: {
        marginLeft: '0.5rem',
    },
    block: {
        display: 'block',
    },
    textSm: {
        fontSize: '0.875rem',
        lineHeight: '1.25rem',
    },
    textGray: {
        color: '#a0aec0',
    },
    button: {
        padding: '10px 20px',
        border: 'none',
        borderRadius: '5px',
        backgroundColor: '#333',
        color: '#fff',
        cursor: 'pointer',
        fontSize: '14px',
        fontWeight: 'bold',
        transition: 'background-color 0.3s ease',
        border: '2px solid #bc8c49',
    },
};