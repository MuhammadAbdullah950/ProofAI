// AlertContext.js
import React, { createContext, useContext, useState } from "react";

// Create the alert context
const AlertContext = createContext();

// Hook to access the alert context
export const useAlert = () => {
    return useContext(AlertContext);
};




// Provider component to wrap the app
export const AlertProvider = ({ children }) => {
    const [alertData, setAlertData] = useState(null);

    // Function to show alert
    const showAlert = (message, type) => {
        setAlertData({ message, type });
    };

    // Function to hide alert
    const hideAlert = () => {
        setAlertData(null);
    };

    return (
        <AlertContext.Provider value={{ showAlert, hideAlert }}>
            {children}


            {alertData && (
                <div
                    style={{
                        position: "fixed",
                        top: "20px",
                        left: "50%",
                        transform: "translateX(-50%)",
                        backgroundColor:
                            alertData.type === "success"
                                ? "#4CAF50"
                                : alertData.type === "error"
                                    ? "#f44336"
                                    : "#2196F3",
                        color: "#fff",
                        padding: "12px 20px",
                        borderRadius: "8px",
                        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
                        zIndex: 9999,
                        fontSize: "16px",
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                        width: "auto",
                        maxWidth: "90%",
                    }}
                >
                    <span style={{ flex: 1 }}>{alertData.message}</span>
                    <button
                        onClick={hideAlert}
                        style={{
                            backgroundColor: "transparent",
                            border: "none",
                            color: "#fff",
                            fontSize: "20px",
                            cursor: "pointer",
                            padding: "0",
                            marginLeft: "10px",
                        }}
                    >
                        &times; {/* Using a multiplication sign for the close button */}
                    </button>
                </div>
            )}

        </AlertContext.Provider>
    );
};
