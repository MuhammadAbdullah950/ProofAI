import React from 'react';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from "../Context/AlertContext";

const useSignup = ({ handleLoginSignup, onPress }) => {

    const proofAiService = useProofAiService();
    const [publicKey, setPublicKey] = React.useState("");
    const [privateKey, setPrivateKey] = React.useState("");
    const { showAlert, hideAlert } = useAlert();


    const handleGenerateKey = async () => {
        const response = await proofAiService.generateKeys();
        hideAlert();
        if (response.error) {
            showAlert("Error in  generating key: " + response.error, "error");
            return;
        }

        setPublicKey(response.pubKey);
        setPrivateKey(response.prvKey);
    }

    const handleLogin = async () => {
        if (publicKey === "" || privateKey === "") {
            showAlert("First generate key! ", "error");
        }
        else {

            const response = await proofAiService.login_using_key(publicKey, privateKey);

            if (response.error) {
                showAlert("Error during loging: " + response.error, "error");
                return;
            }

            if (response.login === "Success") {
                hideAlert();
                onPress();
            } else {
                showAlert("Error during loging: " + response.message, "error");
            }
        }
    }

    return {
        handleGenerateKey,
        handleLogin,
        publicKey,
        privateKey,
    }
};

export default useSignup;
