import React from "react";
import { ArrowRight, Plus, Pickaxe, Settings } from "lucide-react";
import useHome from "../hooks/useHome";

const Home = () => {
    const {
        handleNewTransaction,
        handleMinedTransactions,
        handleMiningTransaction,
        ActionButton
    } = useHome();

    return (
        <div className="min-h-screen flex items-center justify-center p-4 sm:p-6">
            <div className="w-full max-w-lg">
                <div className="backdrop-blur-lg bg-white/10 rounded-2xl shadow-2xl overflow-hidden">
                    <div className="p-6 sm:p-8 border-b border-white/10">
                        <div className="flex items-center justify-between">
                            <h1 className="flex items-center gap-3 text-2xl sm:text-3xl font-bold text-white">
                                <span className="text-3xl sm:text-4xl">💳</span>
                                <span className="bg-gradient-to-r from-amber-200 to-amber-400 bg-clip-text text-transparent">
                                    Transaction Dashboard
                                </span>
                            </h1>
                        </div>
                    </div>

                    <div className="p-6 sm:p-8">
                        <div className="space-y-4">
                            <button
                                onClick={handleNewTransaction}
                                className="w-full group relative overflow-hidden rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 p-px focus:outline-none focus:ring-2 focus:ring-amber-400 focus:ring-offset-2 focus:ring-offset-gray-900"
                            >
                                <div className="relative flex items-center justify-between w-full px-6 py-4 bg-gray-900 rounded-xl transition-all duration-300 group-hover:bg-transparent">
                                    <div className="flex items-center gap-4">
                                        <Plus className="w-6 h-6 text-amber-400" />
                                        <span className="text-lg font-semibold text-white">New Transaction</span>
                                    </div>
                                    <ArrowRight className="w-5 h-5 text-amber-400 transform transition-transform duration-300 group-hover:translate-x-1" />
                                </div>
                            </button>

                            <button
                                onClick={handleMinedTransactions}
                                className="w-full group relative overflow-hidden rounded-xl bg-gradient-to-r from-blue-500 to-blue-600 p-px focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-gray-900"
                            >
                                <div className="relative flex items-center justify-between w-full px-6 py-4 bg-gray-900 rounded-xl transition-all duration-300 group-hover:bg-transparent">
                                    <div className="flex items-center gap-4">
                                        <Pickaxe className="w-6 h-6 text-blue-400" />
                                        <span className="text-lg font-semibold text-white">View Mined Blocks</span>
                                    </div>
                                    <ArrowRight className="w-5 h-5 text-blue-400 transform transition-transform duration-300 group-hover:translate-x-1" />
                                </div>
                            </button>

                            <button
                                onClick={handleMiningTransaction}
                                className="w-full group relative overflow-hidden rounded-xl bg-gradient-to-r from-purple-500 to-purple-600 p-px focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-900"
                            >
                                <div className="relative flex items-center justify-between w-full px-6 py-4 bg-gray-900 rounded-xl transition-all duration-300 group-hover:bg-transparent">
                                    <div className="flex items-center gap-4">
                                        <Settings className="w-6 h-6 text-purple-400 animate-spin-slow" />
                                        <span className="text-lg font-semibold text-white">Currently Mining Block</span>
                                    </div>
                                    <ArrowRight className="w-5 h-5 text-purple-400 transform transition-transform duration-300 group-hover:translate-x-1" />
                                </div>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;