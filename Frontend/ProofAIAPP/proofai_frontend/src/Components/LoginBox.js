import React from 'react';
import { useProofAiService } from '../ProofaiServiceContext';
import { useAlert } from "../Context/AlertContext";

const LoginBox = ({ handleLoginSignup, onPress }) => {
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

  return (
    <div style={styles.maxWidth}>
      <div style={{ ...styles.boxShadow, ...styles.bgGray }}>
        <div style={styles.padding}>
          <h2 style={{ ...styles.textCenter, ...styles.textWhite, ...styles.text3xl, ...styles.fontExtrabold }}>
            Welcome Back
          </h2>
          <p style={{ ...styles.mt4, ...styles.textCenter, ...styles.textGray }}>Sign in to continue</p>

          <div style={styles.roundedShadow}>
            <div>
              <label htmlFor="publicKey" style={styles.label}>Enter Public Key</label>
              <textarea
                placeholder="Public Key"
                style={styles.input}
                required
                autoComplete="publicKey"
                name="publicKey"
                id="publicKey"
                value={publicKey}
                onChange={handleSetPublicKey}
              />
            </div>

            <div>
              <label htmlFor="privateKey" style={styles.label}>Enter Private Key</label>
              <textarea
                placeholder="Private Key"
                style={styles.input}
                required
                autoComplete="privateKey"
                type="password"
                name="privateKey"
                id="privateKey"
                value={privateKey}
                onChange={handleSetPrivateKey}
              />
            </div>
          </div>

          <div style={{ ...styles.flex, ...styles.itemsCenter, ...styles.justifyBetween, ...styles.mt4 }}>
            <div style={styles.flex}>
              <input
                style={styles.checkbox}
                type="checkbox"
                name="remember-me"
                id="remember-me"
                checked={rememberMe}
                onChange={handleRememberme}
              />
              <label style={styles.ml2} htmlFor="remember-me">Remember me</label>
            </div>
          </div>

          <div style={{ ...styles.flex, ...styles.mt4 }}>
            <button style={styles.button} type="button" onClick={handleLogin}>Sign In</button>
            <button style={styles.button} type="button" onClick={handleLoginSignup}>Sign Up</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginBox;

const styles = {
  label: {
    fontSize: '0.975rem',
    lineHeight: '1.25rem',
    color: '#cbd5e0',
    marginBottom: '0.5rem',
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
  roundedShadow: {
    borderRadius: '0.375rem',
    boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
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
    fontSize: '0.875rem',
    lineHeight: '1.25rem',
    resize: 'none',
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
    margin: '5px',
  },
};
