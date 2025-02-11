import React from "react";
import ReactJson from "react-json-view";
import useMinedTrasnaction from "../hooks/useMinedTrasnaction";

const MinedTrasnaction = ({ handleToggleTransactions }) => {
    const { jsonData, handleBack } = useMinedTrasnaction({ handleToggleTransactions });

    return (
        <div className="flex flex-col items-center p-6 bg-gray-200 rounded-xl shadow-lg w-full sm:w-2/3 mx-auto">
            <h1 className="text-lg sm:text-xl font-bold text-gray-800 mb-4">
                Mined Transactions
            </h1>
            <div className="w-full bg-gray-900 p-4 rounded-lg shadow-md overflow-x-auto">
                <ReactJson
                    src={jsonData}
                    theme="monokai"
                    collapsed={false}
                    displayDataTypes={false}
                />
            </div>
            <button
                onClick={handleBack}
                className="mt-4 px-6 py-2 bg-blue-700 hover:bg-blue-800 text-white font-semibold rounded-lg border-2 border-gray-700 transition duration-300"
            >
                Back
            </button>
        </div>
    );
};

export default MinedTrasnaction;
