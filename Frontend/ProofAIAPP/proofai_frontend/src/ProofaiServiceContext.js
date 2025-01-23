
import React, { useContext, createContext } from "react"; // why we use useContext and createContext here? 
import ProofAiService from "./Services/ProofAiService";

const ServiceContext = createContext(null);
export const ServiceProvider = ({ children }) => {
    const proofAiServiceInstance = new ProofAiService(); // Create an instance

    return (
        <ServiceContext.Provider value={proofAiServiceInstance}>{children}</ServiceContext.Provider>
    );
}

export const useProofAiService = () => {
    return useContext(ServiceContext);
}