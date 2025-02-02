
import React from "react";
import ReactJson from 'react-json-view';

const MinedTrasnaction = ({ handleToggleTransactions }) => {
    const jsonData = {
        transactionId: "12345",
        status: "mining",
        details: {
            amount: 100,
            sender: "Alice",
            receiver: "Bob"
        }
    };

    const handleBack = () => {
        handleToggleTransactions()
    }

    return (
        <div >
            <h1 style={styles.label} >Mined Transactions</h1>
            <ReactJson src={jsonData} theme="monokai" collapsed={false} displayDataTypes={false} />
            <button style={styles.button} onClick={handleBack} >Back</button>
        </div>
    );
}

export default MinedTrasnaction;


const styles = {
    label: {
        fontSize: '27px',
        fontWeight: 'bold',
        color: '#333',
        textAlign: 'center',
    },

    box: {
        display: 'flex',
        flexDirection: 'column',
        gap: '10px',
        padding: '20px',
        backgroundColor: "rgb(152, 156, 158)",
        borderRadius: '10px',
        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
        width: '50%',
    },

    input: {
        display: 'block',
        width: '100%',
        padding: '8px',
        marginTop: '5px',
        borderRadius: '5px',
        border: '1px solid #ccc',
        boxSizing: 'border-box',
    },

    buttonContianer: {
        gap: '10px',
        display: 'flex',
        justifyContent: 'space-around',
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
    buttonHover: {
        backgroundColor: '#0056b3',
    },
}