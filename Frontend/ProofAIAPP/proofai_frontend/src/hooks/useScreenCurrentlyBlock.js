import React from "react";
import { FaDownload } from 'react-icons/fa';
import ReactJson from 'react-json-view';
import { useProofAiService } from "../ProofaiServiceContext";
import { useAlert } from "../Context/AlertContext";
import { useNavigate } from "react-router-dom";

const useScreenCurrentlyBlock = () => {

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


    return {
        showTransaction,
        setShowTransaction,
        block,
        setBlock,
        ProofAiService,
        handleBack,
        getBlock
    }
};

export default useScreenCurrentlyBlock;
