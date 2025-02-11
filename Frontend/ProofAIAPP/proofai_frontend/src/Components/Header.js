import React from "react";
import { LogOut, Settings } from "lucide-react";
import { useProofAiService } from "../ProofaiServiceContext";
import useHeader from "../hooks/useHeader";

const Header = ({ handleLogout }) => {
    const { handleRoleChange, role } = useHeader();

    return (
        <header className="fixed top-0 left-0 w-full bg-gray-800 shadow-md z-50 px-4 py-3 flex items-center justify-between">
            <div className="text-2xl font-bold text-white">ProofAI</div>

            <div className="flex items-center space-x-4">
                <select
                    onChange={handleRoleChange}
                    className="bg-gray-700 text-white px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                    <option value="Miner">Miner</option>
                    <option value="Validator">Validator</option>
                </select>

                <button
                    onClick={handleLogout}
                    className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 transition-colors flex items-center space-x-2"
                >
                    <LogOut size={18} />
                    <span>Logout</span>
                </button>
            </div>
        </header>
    );
}

export default Header;