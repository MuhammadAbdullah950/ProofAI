import React from "react";

const Footer = () => {
    return (
        <div style={styles.container} >
            <div  style={styles.Footer} >Copyright © 2024 CampusView. All rights reserved.</div>
        </div>
    );
}
export default Footer;

const styles = ({
    container: {
        display: 'flex',
        justifyContent: 'center', 
        alignItems: 'center', 
        width: '100%',
        padding: '10px', 
        position: 'absolute', 
        bottom: 0,
        boxShadow: '0 -2px 5px rgba(0, 0, 0, 0.1)', 
        backgroundColor: '#2c3e50',
    },
    footer: {
        color: 'white',
        textAlign: 'center',
        padding: '10px',
        fontSize: '14px',
    }
})