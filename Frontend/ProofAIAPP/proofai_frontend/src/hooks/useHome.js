import React, { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useProofAiService } from "../ProofaiServiceContext";
import { ArrowRight, Plus, Pickaxe, Settings } from "lucide-react";


const useHome = () => {
    const navigate = useNavigate()

    const handleNewTransaction = () => {
        navigate('/NewTransaction')
    }

    const handleMinedTransactions = () => {
        navigate('/MinedTransactions')
    }

    const ProofAiService = useProofAiService()

    const handleMiningTransaction = async () => {
        navigate('/CurrentlyBlock')
    }


    const ActionButton = ({ onClick, label, icon: Icon, variant }) => {
        const baseStyles = "w-full p-4 rounded-lg font-medium transition-all duration-300 flex items-center justify-between group";
        const variants = {
            primary: "bg-indigo-600 hover:bg-indigo-700 text-white",
            secondary: "bg-teal-600 hover:bg-teal-700 text-white",
            tertiary: "bg-purple-600 hover:bg-purple-700 text-white"
        };

        return (
            <button
                onClick={onClick}
                className={`${baseStyles} ${variants[variant]}`}
            >
                <span className="flex items-center gap-3">
                    <Icon className="w-5 h-5" />
                    {label}
                </span>
                <ArrowRight className="w-4 h-4 transition-transform duration-300 group-hover:translate-x-1" />
            </button>
        );
    };



    return {
        handleNewTransaction,
        handleMinedTransactions,
        handleMiningTransaction,
        ActionButton,

    }
}

export default useHome;

