
import React, { useContext, createContext } from "react";
import ProofAiService from "./Services/ProofAiService";

const ServiceContext = createContext(null);
export const ServiceProvider = ({ children }) => {
    const proofAiServiceInstance = new ProofAiService();
    return (
        <ServiceContext.Provider value={proofAiServiceInstance}>{children}</ServiceContext.Provider>
    );
}
export const useProofAiService = () => {
    return useContext(ServiceContext);
}