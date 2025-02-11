import React from "react";
import ReactJson from "react-json-view";
import { ArrowLeft } from "lucide-react";
import useScreenCurrentlyBlock from "../hooks/useScreenCurrentlyBlock";

const CurrentlyBlock = () => {
    const { block, handleBack } = useScreenCurrentlyBlock();

    return (
        <div className="min-h-screen  p-6">
            <div className="mx-auto max-w-7xl space-y-6">

                <div className="rounded-lg bg-white/5 p-6 backdrop-blur-lg">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <div className="rounded-lg bg-slate-700 p-2">
                                <svg
                                    className="h-6 w-6 text-emerald-400"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                >
                                    <rect x="2" y="2" width="20" height="20" rx="2" />
                                    <path d="M12 2v20" />
                                    <path d="M2 12h20" />
                                </svg>
                            </div>
                            <h1 className="text-3xl font-bold text-white">
                                Currently Mining Block
                            </h1>
                        </div>
                    </div>
                </div>


                <div className="overflow-hidden rounded-lg bg-slate-800/50 backdrop-blur-lg">
                    <div className="border-b border-slate-700 bg-slate-800 p-4">
                        <h2 className="text-sm font-medium text-slate-300">
                            Block Details
                        </h2>
                    </div>
                    <div className="p-6">
                        <div className="rounded-lg bg-slate-900">
                            <ReactJson
                                src={block}
                                theme="monokai"
                                collapsed={false}
                                displayDataTypes={false}
                                displayObjectSize={true}
                                style={{
                                    backgroundColor: 'transparent',
                                    fontSize: '14px',
                                    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
                                    padding: '1rem'
                                }}
                                iconStyle="square"
                                enableClipboard={true}
                                collapseStringsAfterLength={80}
                            />
                        </div>
                    </div>
                </div>


                <button
                    onClick={handleBack}
                    className="flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-slate-800"
                >
                    <ArrowLeft className="h-4 w-4" />
                    Back
                </button>
            </div>
        </div>
    );
};

export default CurrentlyBlock;