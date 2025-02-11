import React from 'react';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from "../Context/AlertContext";

const useLogin = ({ handleLoginSignup, onPress }) => {
    const { showAlert, hideAlert } = useAlert();
    const proofAiService = useProofAiService();
    const [publicKey, setPublicKey] = React.useState("");
    const [privateKey, setPrivateKey] = React.useState("");
    const [rememberMe, setRememberMe] = React.useState(false);

    React.useEffect(() => {
        const savedPublicKey = localStorage.getItem("publicKey");
        const savedPrivateKey = localStorage.getItem("privateKey");

        if (savedPublicKey && savedPrivateKey) {
            setPublicKey(savedPublicKey);
            setPrivateKey(savedPrivateKey);
            setRememberMe(true);
        }
    }, []);

    const handleLogin = async () => {
        const response = await proofAiService.login_using_key(publicKey, privateKey);

        if (response.error) {
            showAlert("Error during login: " + response.error, "error");
            return;
        }

        if (response.login === "Success") {
            if (rememberMe) {
                localStorage.setItem("publicKey", publicKey);
                localStorage.setItem("privateKey", privateKey);
            } else {
                localStorage.removeItem("publicKey");
                localStorage.removeItem("privateKey");
            }

            hideAlert();
            onPress();
        } else {
            showAlert("Error during login: " + response.message, "error");
        }
    };

    const handleSetPublicKey = (event) => {
        setPublicKey(event.target.value);
        event.target.style.height = "auto";
        event.target.style.height = event.target.scrollHeight + "px";
    };

    const handleSetPrivateKey = (event) => {
        setPrivateKey(event.target.value);
        event.target.style.height = "auto";
        event.target.style.height = event.target.scrollHeight + "px";
    };

    const handleRememberme = (event) => {
        setRememberMe(event.target.checked);
    };

    return {
        handleLogin,
        handleSetPublicKey,
        handleSetPrivateKey,
        handleRememberme,
        rememberMe,
        publicKey,
        privateKey
    }
};

export default useLogin;
