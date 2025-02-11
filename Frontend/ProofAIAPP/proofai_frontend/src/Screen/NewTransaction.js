import React from 'react';
import { ClipLoader } from 'react-spinners';
import { ArrowLeft, Upload, Database, FileCode, HardDrive } from 'lucide-react';
import IPFS from '../Components/Ipfs';
import MiningTransaction from '../Components/MiningTransaction';
import useNewTransaction from '../hooks/useNewTransaction';

const NewTransaction = () => {
    const {
        ipfsCont,
        setipfsCont,
        mineTransaction,
        setmineTransaction,
        setmodelCID,
        setDatasetCID,
        responseWait,
        transactionJson,
        handleMine,
        handleIpfs,
        handleUserInstructions,
        handleBackButton,
    } = useNewTransaction();

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-700 flex items-center justify-center p-4 sm:p-6">

            {ipfsCont && <IPFS setipfsCont={setipfsCont} />}


            {mineTransaction && (
                <MiningTransaction
                    setmineTransaction={setmineTransaction}
                    transactionJson={transactionJson}
                />
            )}

            {responseWait && (
                <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex justify-center items-center">
                    <div className="bg-gray-800 p-8 rounded-2xl shadow-xl flex flex-col items-center">
                        <ClipLoader color="#F59E0B" loading={true} size={50} />
                        <p className="text-white mt-4">Processing Transaction...</p>
                    </div>
                </div>
            )}


            {(!ipfsCont && !mineTransaction) && (
                <div className="w-full max-w-lg">
                    <div className="backdrop-blur-lg bg-white/10 rounded-2xl shadow-2xl overflow-hidden">

                        <div className="p-6 sm:p-8 border-b border-white/10 flex items-center justify-between">
                            <h1 className="text-2xl sm:text-3xl font-bold text-white flex items-center gap-3">
                                <FileCode className="w-8 h-8 text-amber-400" />
                                <span className="bg-gradient-to-r from-amber-200 to-amber-400 bg-clip-text text-transparent">
                                    New Transaction
                                </span>
                            </h1>
                            <button
                                onClick={handleBackButton}
                                className="p-2 rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white transition-colors"
                            >
                                <ArrowLeft className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="p-6 sm:p-8 space-y-6">

                            <div className="space-y-2">
                                <label className="flex items-center gap-2 text-sm font-medium text-gray-300">
                                    <Database className="w-4 h-4 text-amber-400" />
                                    Model CID
                                </label>
                                <input
                                    className="w-full px-4 py-3 bg-gray-800 text-white border border-gray-700 rounded-xl focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors placeholder-gray-500"
                                    placeholder="Enter the CID of your model"
                                    onChange={(event) => setmodelCID(event.target.value)}
                                />
                            </div>

                            <div className="space-y-2">
                                <label className="flex items-center gap-2 text-sm font-medium text-gray-300">
                                    <HardDrive className="w-4 h-4 text-amber-400" />
                                    Dataset CID
                                </label>
                                <input
                                    className="w-full px-4 py-3 bg-gray-800 text-white border border-gray-700 rounded-xl focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors placeholder-gray-500"
                                    placeholder="Enter the CID of your dataset"
                                    onChange={(event) => setDatasetCID(event.target.value)}
                                />
                            </div>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                <button
                                    onClick={handleMine}
                                    className="flex items-center justify-center gap-2 py-3 px-4 bg-gradient-to-r from-amber-500 to-amber-600 text-white font-semibold rounded-xl hover:from-amber-600 hover:to-amber-700 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 focus:ring-offset-gray-900"
                                >
                                    <Database className="w-5 h-5" />
                                    Mine Transaction
                                </button>
                                <button
                                    onClick={handleIpfs}
                                    className="flex items-center justify-center gap-2 py-3 px-4 bg-gradient-to-r from-blue-500 to-blue-600 text-white font-semibold rounded-xl hover:from-blue-600 hover:to-blue-700 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900"
                                >
                                    <Upload className="w-5 h-5" />
                                    Upload to IPFS
                                </button>
                            </div>

                            <button
                                onClick={handleUserInstructions}
                                className="w-full py-3 px-4 bg-gray-800 text-gray-300 font-medium rounded-xl border border-gray-700 hover:bg-gray-700 hover:text-white transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 focus:ring-offset-gray-900"
                            >
                                View Transaction Instructions
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default NewTransaction;