import React from 'react';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from "../Context/AlertContext";

const SignupBox = ({ handleLoginSignup, onPress }) => {

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

    return (
        <div style={styles.maxWidth}>
            <div style={{ ...styles.boxShadow, ...styles.bgGray }}>
                <div style={styles.padding}>
                    <h2 style={{ ...styles.textCenter, ...styles.textWhite, ...styles.text3xl, ...styles.fontExtrabold }}> Welcome To ProofAI</h2>
                    <div style={styles.roundedShadow}>
                        <div>
                            <label htmlFor="publicKey" style={{ ...styles.label, fontWeight: "bold", color: "#c8deaf" }} >Public Key</label>
                            <textarea
                                style={{ ...styles.input, overflowY: "auto", resize: "none", height: "auto" }}
                                required
                                name="publicKey"
                                id="publicKey"
                                value={publicKey}
                                readOnly
                            />
                        </div>

                        <div style={styles.mt4}>
                            <label htmlFor="privateKey" style={{ ...styles.label, fontWeight: "bold", color: "#ccdeaf" }} >Private Key</label>
                            <textarea
                                style={{ ...styles.input, overflowY: "auto", resize: "none", height: "auto" }}
                                name="privateKey"
                                id="privateKey"
                                value={privateKey}
                                readOnly
                            />
                        </div>
                    </div>
                    <div style={{ display: "flex", flexDirection: "row", gap: '20px', marginTop: '20px' }}>
                        <button style={styles.button} onClick={handleGenerateKey}>Generate Key</button>
                        <button style={styles.button} onClick={handleLogin}  >Login-in</button>
                    </div>
                </div>
            </div>
        </div >
    );
};

export default SignupBox;

const styles = {
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
        outline: 'none',
        zIndex: '10',
        fontSize: '0.875rem',
        lineHeight: '1.25rem',
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