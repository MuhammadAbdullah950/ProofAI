
import React from "react";
import ReactJson from 'react-json-view';

const useMinedTrasnaction = ({ handleToggleTransactions }) => {
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

    return {
        jsonData,
        handleBack
    }
}

export default useMinedTrasnaction;

