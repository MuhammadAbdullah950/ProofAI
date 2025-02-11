import React from 'react';
import useLogin from '../hooks/useLogin';

const LoginBox = ({ handleLoginSignup, onPress }) => {
  const {
    handleSetPublicKey,
    handleSetPrivateKey,
    handleRememberme,
    handleLogin,
    publicKey,
    privateKey,
    rememberMe,
  } = useLogin({ handleLoginSignup, onPress });

  return (
    <div className="w-full">
      <div className="bg-gray-800 shadow-2xl rounded-lg">
        <div className="p-4 sm:p-6 md:p-8">
          <div className="">
            <div className="flex-1">
              <label htmlFor="publicKey" className="block text-sm text-gray-300 mb-2">
                Enter Public Key
              </label>
              <textarea
                placeholder="Public Key"
                className="w-full min-h-[80px] max-h-[120px] px-3 py-2 bg-gray-700 text-white border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 resize-y"
                required
                autoComplete="publicKey"
                name="publicKey"
                id="publicKey"
                value={publicKey}
                onChange={handleSetPublicKey}
              />
            </div>

            <div className="flex-1">
              <label htmlFor="privateKey" className="block text-sm text-gray-300 mb-2">
                Enter Private Key
              </label>
              <textarea
                placeholder="Private Key"
                className="w-full min-h-[80px] max-h-[120px] px-3 py-2 bg-gray-700 text-white border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 resize-y"
                required
                autoComplete="privateKey"
                type="password"
                name="privateKey"
                id="privateKey"
                value={privateKey}
                onChange={handleSetPrivateKey}
              />
            </div>

            <div className="flex items-center">
              <input
                type="checkbox"
                name="remember-me"
                id="remember-me"
                className="h-4 w-4 text-amber-600 focus:ring-amber-500 border-gray-300 rounded"
                checked={rememberMe}
                onChange={handleRememberme}
              />
              <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-300">
                Remember me
              </label>
            </div>

            <div className="flex flex-col sm:flex-row gap-4 mt-2">
              <button
                className="flex-1 py-2 px-4 bg-gray-700 text-white font-bold rounded-md border-2 border-amber-600 hover:bg-amber-700 transition duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-amber-500"
                type="button"
                onClick={handleLogin}
              >
                Sign In
              </button>
              <button
                className="flex-1 py-2 px-4 bg-gray-700 text-white font-bold rounded-md border-2 border-amber-600 hover:bg-amber-700 transition duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-amber-500"
                type="button"
                onClick={handleLoginSignup}
              >
                Sign Up
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginBox;