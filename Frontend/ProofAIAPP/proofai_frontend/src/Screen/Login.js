import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import LoginBox from "../Components/LoginBox";
import Footer from "../Components/Footer";
import SignupBox from "../Components/SignupBox";
import useScreenlogin from "../hooks/useScreenlogin";

const Login = () => {
    const { isLogin, handleLogin, handleLoginSignup, boxHeight } = useScreenlogin();



    return (
        <div className="min-h-screen  flex flex-col">
            <div
                className=" flex items-center justify-center  "
                style={{ minHeight: boxHeight }}
            >
                <div className="w-full max-w-md">
                    <div className=" rounded-xl shadow-2xl p-4 sm:p-4 flex flex-col h-full border-4 border-gray-300">
                        <div className="text-center mb-6">
                            <h2 className="text-3xl font-bold text-white">
                                {isLogin ? "Welcome Back" : "Create Account"}
                            </h2>
                            <p className="mt-2 text-sm text-gray-600">
                                {isLogin
                                    ? "Sign in to access your account"
                                    : "Join us to get started"}
                            </p>
                        </div>

                        <div className="flex-1 flex flex-col justify-center">
                            {isLogin ? (
                                <LoginBox
                                    onPress={handleLogin}
                                    handleLoginSignup={handleLoginSignup}
                                />
                            ) : (
                                <SignupBox
                                    handleLoginSignup={handleLoginSignup}
                                    onPress={handleLogin}
                                />
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login;