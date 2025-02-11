import React from 'react';
import { ClipLoader } from 'react-spinners';
import { Upload, Copy, ArrowLeft } from 'lucide-react';
import useIpfs from '../hooks/useIpfs';

const IPFS = ({ setipfsCont }) => {
    const { uploadResponse, waitingIcon, handleUpload, handleCopyCID, handleBackButton } = useIpfs({ setipfsCont });

    return (
        <div className="min-h-screen flex items-center justify-center pb-16">
            <div className="mx-auto max-w-2xl">
                <div className="overflow-hidden rounded-xl bg-white/10 backdrop-blur-lg">

                    <div className="border-b border-white/10 bg-white/5 p-6">
                        <h2 className="text-center text-2xl font-bold text-white">
                            Upload Folder to IPFS
                        </h2>
                    </div>


                    <div className="space-y-6 p-6">

                        <div className="rounded-lg border-2 border-dashed border-slate-400/25 p-6 text-center">
                            <input
                                type="file"
                                id="file"
                                className="hidden"
                                webkitdirectory="true"
                                directory="true"
                                multiple
                                onChange={handleUpload}
                            />
                            <label
                                htmlFor="file"
                                className="group cursor-pointer"
                            >
                                <div className="space-y-4">
                                    <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-slate-700 group-hover:bg-slate-600">
                                        <Upload className="h-6 w-6 text-slate-200" />
                                    </div>
                                    <div>
                                        <p className="text-sm font-medium text-slate-300">
                                            Drag and drop your folder here, or click to browse
                                        </p>
                                        <p className="mt-1 text-xs text-slate-400">
                                            All files within the folder will be uploaded
                                        </p>
                                    </div>
                                </div>
                            </label>
                        </div>


                        {waitingIcon && (
                            <div className="flex justify-center py-4">
                                <ClipLoader color="#94a3b8" loading={true} size={40} />
                            </div>
                        )}


                        {uploadResponse && (
                            <div className="overflow-hidden rounded-lg bg-slate-800/50 p-4">
                                <p className="mb-2 text-sm font-medium text-slate-300">
                                    CID of Uploaded Folder:
                                </p>
                                <div className="flex items-center gap-2">
                                    <code className="flex-1 overflow-x-auto rounded bg-slate-900/50 p-2 text-sm text-slate-300">
                                        {uploadResponse}
                                    </code>
                                    <button
                                        onClick={handleCopyCID}
                                        className="flex items-center gap-2 rounded-lg bg-emerald-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-slate-800"
                                    >
                                        <Copy className="h-4 w-4" />
                                        Copy
                                    </button>
                                </div>
                            </div>
                        )}


                        <button
                            onClick={handleBackButton}
                            className="flex items-center gap-2 rounded-lg bg-slate-700 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 focus:ring-offset-slate-800"
                        >
                            <ArrowLeft className="h-4 w-4" />
                            Back
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default IPFS;