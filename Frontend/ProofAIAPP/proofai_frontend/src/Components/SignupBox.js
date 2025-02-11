import React from 'react';
import useSignup from '../hooks/useSignup';

const SignupBox = ({ handleLoginSignup, onPress }) => {
    const {
        handleGenerateKey,
        handleLogin,
        publicKey,
        privateKey
    } = useSignup({ handleLoginSignup, onPress });

    return (
        <div className="w-full">
            <div className="bg-gray-800 shadow-2xl rounded-lg">
                <div className="p-4 sm:p-6 md:p-8">
                    <div className="">
                        <div className="flex-1">
                            <label
                                htmlFor="publicKey"
                                className="block text-sm font-medium text-green-300 mb-2"
                            >
                                Public Key
                            </label>
                            <div className="relative">
                                <textarea
                                    className="w-full min-h-[80px] max-h-[120px] px-3 py-2 bg-gray-700 text-white border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
                                    required
                                    name="publicKey"
                                    id="publicKey"
                                    value={publicKey}
                                    readOnly
                                />
                            </div>
                        </div>

                        <div className="flex-1">
                            <label
                                htmlFor="privateKey"
                                className="block text-sm font-medium text-green-300 mb-2"
                            >
                                Private Key
                            </label>
                            <div className="relative">
                                <textarea
                                    className="w-full min-h-[80px] max-h-[120px] px-3 py-2 bg-gray-700 text-white border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
                                    name="privateKey"
                                    id="privateKey"
                                    value={privateKey}
                                    readOnly
                                />
                            </div>
                        </div>

                        <div className="flex flex-col sm:flex-row gap-4">
                            <button
                                className="flex-1 py-2 px-4 bg-gray-700 text-white font-bold rounded-md border-2 border-amber-600 hover:bg-amber-700 transition duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50 disabled:cursor-not-allowed"
                                onClick={handleGenerateKey}
                                disabled={!!publicKey}
                            >
                                Generate Key
                            </button>
                            <button
                                className="flex-1 py-2 px-4 bg-gray-700 text-white font-bold rounded-md border-2 border-amber-600 hover:bg-amber-700 transition duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-amber-500"
                                onClick={handleLogin}
                            >
                                Log In
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SignupBox;