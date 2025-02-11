import React from "react";
import ReactJson from 'react-json-view';
import { ClipLoader } from 'react-spinners';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from '../Context/AlertContext';
import { FaCheckCircle } from 'react-icons/fa';
import { useEffect } from "react";

const useMiningTransaction = ({ setmineTransaction, transactionJson }) => {
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

    return {
        isTransactionMined,
        chkTransactionStatus,
        handleBack,

    }
};

export default useMiningTransaction;
