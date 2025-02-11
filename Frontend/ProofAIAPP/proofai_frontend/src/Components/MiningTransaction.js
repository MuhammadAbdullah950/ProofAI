import React from "react";
import ReactJson from 'react-json-view';
import { ClipLoader } from 'react-spinners';
import { FaCheckCircle } from 'react-icons/fa';
import useMiningTransaction from "../hooks/useMiningTransaction";

const MiningTransaction = ({ setmineTransaction, transactionJson }) => {
    const { isTransactionMined, chkTransactionStatus, handleBack } = useMiningTransaction({ setmineTransaction, transactionJson });

    return (
        <div className="flex flex-col items-center gap-4 p-6 bg-gray-300 rounded-xl shadow-lg max-w-3xl w-full mx-auto mt-16">
            <h1 className="text-2xl font-bold text-gray-800 text-center">Mining Transaction</h1>

            {!isTransactionMined && (
                <>
                    <div className="flex justify-center items-center h-12">
                        <ClipLoader color={'#123abc'} loading={true} size={50} />
                    </div>

                    <div className="w-full h-72 p-4 bg-gray-900 text-white rounded-md shadow-md overflow-auto">
                        <ReactJson
                            src={transactionJson}
                            theme="monokai"
                            collapsed={false}
                            displayDataTypes={false}
                            displayObjectSize={true}
                        />
                    </div>
                </>
            )}

            {isTransactionMined && <FaCheckCircle className="text-green-500 text-2xl" />}

            <div className="flex gap-4 mt-4">
                <button
                    className="px-4 py-2 bg-blue-700 text-white font-semibold rounded-md border-2 border-gray-700 hover:bg-blue-800 transition"
                    onClick={chkTransactionStatus}
                >
                    Check Status
                </button>
                <button
                    className="px-4 py-2 bg-blue-700 text-white font-semibold rounded-md border-2 border-gray-700 hover:bg-blue-800 transition"
                    onClick={handleBack}
                >
                    Back
                </button>
            </div>
        </div>
    );
};

export default MiningTransaction;
