import React, { useState } from 'react';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from "../Context/AlertContext";
import { ClipLoader } from 'react-spinners';
import axios from 'axios';

const useIpfs = ({ setipfsCont }) => {
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

    return {
        uploadResponse,
        ipfsContent,
        waitingIcon,
        handleUpload,
        handleCopyCID,
        handleBackButton
    }
};

export default useIpfs;
