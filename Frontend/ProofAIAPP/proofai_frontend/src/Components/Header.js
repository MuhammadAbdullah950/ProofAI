import React from "react";
import { useProofAiService } from "../ProofaiServiceContext";

const Header = ({ handleLogout }) => {
    const ProofAiService = useProofAiService();
    const [role, setRole] = React.useState(true);


    const handleRoleChange = async (e) => {
        const minerRole = e.target.value;
        const response = await ProofAiService.setRole(minerRole);
        alert(response);
        if (response.error) {
            alert(response.error);
            setRole(role);
            return;
        }
        setRole(minerRole);
    }


    return (
        <div style={styles.headerContainer} >
            <label style={styles.appName} >ProofAI</label>

            <div style={styles.buttonContainer}>
                <select style={styles.Button} onChange={handleRoleChange} >
                    <option value="Miner" selected>Miner</option>
                    <option value="Validator" >Validator</option>
                </select>
                <button style={styles.Button} onClick={handleLogout} >Logout</button>
            </div>

        </div>
    );
}
export default Header;

const styles = {
    headerContainer: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '10px 20px',
        backgroundColor: '#2c3e50',
        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
        boxSizing: 'border-box',
    },
    appName: {
        fontSize: '24px',
        fontWeight: 'bold',
        color: '#ffffff',
    },

    buttonContainer: {
        display: 'flex',
        gap: '10px',
    },

    Button: {
        padding: '8px 16px',
        border: 'none',
        borderRadius: '5px',
        backgroundColor: '#333',
        color: '#fff',
        cursor: 'pointer',
        fontSize: '14px',
        transition: 'background-color 0.3s ease',
        border: '2px solid #595f66',
    },
    ButtonHover: { // Button Hover
        backgroundColor: '#0056b3',
    },
    ButtonFocus: {
        outline: 'none',
        boxShadow: '0 0 0 2px rgba(0, 123, 255, 0.5)',
    },

};
