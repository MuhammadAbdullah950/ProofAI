import React from "react";
import { useNavigate } from "react-router-dom";
import { ChevronLeft, Box, Filter } from "lucide-react";
import useScreenMinedTransactions from "../hooks/useScreenMinedTransactions";

const MinedTransactions = () => {
    const {
        minedBlock,
        handleBack,
        handleNavigate,
        handleFilterBlocks,
    } = useScreenMinedTransactions();

    return (
        <div className="min-h-screen  p-4 sm:p-6">

            <div className="max-w-6xl mx-auto">
                <div className="mb-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <button
                            onClick={handleBack}
                            className="p-2 rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-amber-500"
                        >
                            <ChevronLeft className="w-5 h-5" />
                        </button>
                        <h1 className="text-2xl sm:text-3xl font-bold text-white">Mined Blocks</h1>
                    </div>

                    <div className="relative">
                        <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                            <Filter className="h-4 w-4 text-gray-400" />
                        </div>
                        <select
                            onChange={handleFilterBlocks}
                            className="pl-10 pr-4 py-2 bg-gray-800 text-white rounded-lg border border-gray-700 focus:border-amber-500 focus:ring-2 focus:ring-amber-500 outline-none appearance-none hover:bg-gray-700 transition-colors"
                        >
                            <option value="All Transactions">All Transactions</option>
                            <option value="Own Transactions">Own Transactions</option>
                        </select>
                    </div>
                </div>


                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {minedBlock && minedBlock.length > 0 ? (
                        minedBlock.map((block) => (
                            <div
                                key={block.BlockNum}
                                className="bg-gray-800/50 backdrop-blur-sm rounded-xl overflow-hidden border border-gray-700 hover:border-amber-500/50 transition-colors"
                            >
                                <div className="p-4 border-b border-gray-700">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-2">
                                            <Box className="w-5 h-5 text-amber-500" />
                                            <span className="text-lg font-semibold text-white">
                                                Block {block.blockNum}
                                            </span>
                                        </div>
                                    </div>
                                </div>

                                <div className="p-4">
                                    <select
                                        onChange={(e) => handleNavigate(e.target.value)}
                                        className="w-full bg-gray-700 text-white rounded-lg p-2 border border-gray-600 focus:border-amber-500 focus:ring-2 focus:ring-amber-500 outline-none"
                                    >
                                        <option value="">Select Transaction</option>
                                        <option value={JSON.stringify(block)}>Complete Block Data</option>
                                        {block.transactions && block.transactions.length > 0 ? (
                                            block.transactions.map((transaction) => (
                                                <option
                                                    key={transaction.transactionId}
                                                    value={JSON.stringify(transaction)}
                                                    className="py-2"
                                                >
                                                    Nonce: {transaction.nonce} | From: {transaction.from?.substring(0, 20)}...
                                                </option>
                                            ))
                                        ) : (
                                            <option disabled>No transactions</option>
                                        )}
                                    </select>
                                </div>
                            </div>
                        ))
                    ) : (
                        <div className="col-span-full flex flex-col items-center justify-center p-8 bg-gray-800/50 backdrop-blur-sm rounded-xl border border-gray-700">
                            <Box className="w-12 h-12 text-gray-600 mb-4" />
                            <p className="text-gray-400 text-lg">No mined blocks found</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default MinedTransactions;