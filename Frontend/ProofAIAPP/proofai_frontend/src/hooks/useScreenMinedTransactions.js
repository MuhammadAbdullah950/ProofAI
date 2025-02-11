import React from "react";
import { useNavigate } from "react-router-dom";
import { useProofAiService } from "../ProofaiServiceContext";
import { useAlert } from "../Context/AlertContext";


const useScreenMinedTransactions = () => {
    const { showAlert, hideAlert } = useAlert();
    const navigate = useNavigate();
    const ProofAiService = useProofAiService();
    const [minedBlock, setMinedBlock] = React.useState([]);
    const [tempMinedBlock, setTempMinedBlock] = React.useState([]);
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
                    setMinedBlock(response.blocks);
                    setTempMinedBlock(response.blocks);
                } else {
                    setMinedBlock([]);
                    showAlert("No mined blocks found.", "success");
                }
            }
        } catch (error) {
            showAlert("Failed to fetch mined blocks.", "error");
            setMinedBlock([]);
        } finally {
            setLoading(false);
        }
    };


    React.useEffect(() => {
        GetBlocks();
    }, []);

    const handleBack = () => {
        hideAlert();
        window.history.back();
    };

    const handleNavigate = (transactionString) => {
        const transaction = JSON.parse(transactionString);
        navigate(`/MinedTransaction`, { state: { transaction } });
    };

    const handleFilterBlocks = async (e) => {

        if (e.target.value !== "Own Transactions") {
            setMinedBlock(tempMinedBlock);
        } else {
            const response = await ProofAiService.getPublicKey();

            if (response.error) {
                showAlert(response.error, "error");
            } else {
                const publicKey = response.pubKey;
                const filteredBlocks = minedBlock
                    .map((block) => {
                        const filteredTransactions = block.transactions?.filter(
                            (transaction) => transaction.from === publicKey
                        );
                        return { ...block, transactions: filteredTransactions };
                    })
                    .filter((block) => block.transactions && block.transactions.length > 0);

                setMinedBlock(filteredBlocks);
            }
        }
    };

    return {
        minedBlock,
        loading,
        handleBack,
        handleNavigate,
        handleFilterBlocks,
        filterValue,
        setFilterValue
    }
};

export default useScreenMinedTransactions;
