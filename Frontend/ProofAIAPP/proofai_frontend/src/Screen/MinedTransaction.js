import React from "react";
import { FaDownload } from "react-icons/fa";
import ReactJson from "react-json-view";
import useScreenMinedTransaction from "../hooks/useScreenMinedTransaction";

const MinedTransaction = () => {
    const {
        transaction,
        handleBack,
        handleModelDownload,
        handlTransactionLogDownload,
        isDiabled,
    } = useScreenMinedTransaction();

    return (
        <div className="flex flex-col min-h-screen p-6 bg-gray-900 text-white">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">Mined Transactions</h1>
                <div className="flex gap-4">
                    <button
                        className={`flex items-center gap-2 px-4 py-2 border border-yellow-400 rounded-md bg-white text-gray-900 text-sm font-medium ${isDiabled ? "opacity-50 cursor-not-allowed" : "hover:bg-yellow-300"
                            }`}
                        onClick={!isDiabled ? handleModelDownload : undefined}
                        disabled={isDiabled}
                    >
                        Download Trained Model <FaDownload />
                    </button>
                    <button
                        className={`flex items-center gap-2 px-4 py-2 border border-yellow-400 rounded-md bg-white text-gray-900 text-sm font-medium ${isDiabled ? "opacity-50 cursor-not-allowed" : "hover:bg-yellow-300"
                            }`}
                        onClick={!isDiabled ? handlTransactionLogDownload : undefined}
                        disabled={isDiabled}
                    >
                        Download Transaction Logs <FaDownload />
                    </button>
                </div>
            </div>

            <div className="bg-gray-700 p-4 rounded-lg shadow-md overflow-auto max-h-[400px]">
                <ReactJson src={transaction} theme="monokai" collapsed={false} displayDataTypes={false} />
            </div>

            <button
                className="mt-4 w-32 px-4 py-2 border border-white text-black bg-yellow-500 rounded-md font-bold hover:bg-yellow-600 transition-all"
                onClick={handleBack}
            >
                Back
            </button>
        </div>
    );
};

export default MinedTransaction;
