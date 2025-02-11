import React, { useEffect, useState } from 'react';
import Login from '../Screen/Login';
import Home from '../Screen/Home';
import Header from '../Components/Header';
import { Routes, Route, useNavigate } from 'react-router-dom';
import NewTransaction from '../Screen/NewTransaction';
import MinedTransactions from '../Screen/MinedTransactions';
import MinedTransaction from '../Screen/MinedTransaction';
import CurrentlyBlock from "../Screen/CurrentlyBlock";
import { useProofAiService } from "../ProofaiServiceContext";
import { useAlert } from "../Context/AlertContext";
import path from 'path-browserify';
import UserInstructions from '../Screen/UserInstructions';

const { ipcRenderer } = window.require("electron");


const useApp = () => {

    const navigate = useNavigate();
    const ProofAiService = useProofAiService();
    const { showAlert, hideAlert } = useAlert();

    const [isServiceAddressSet, setIsServiceAddressSet] = useState(false);
    const [serviceMachineAddr, setServiceMachineAddr] = useState("");
    const [showModal, setShowModal] = useState(false);
    const [serviceButtondisable, setServiceButtondisable] = useState(false);

    const handleLogout = async () => {
        hideAlert();
        const response = await ProofAiService.getCurrentlyMinBlock();
        if (response.error) {
            showAlert(response.error, "error");
            return;
        }

        if (response.block !== "null") {
            const userConfirmation = window.confirm("Do you want to continue with the current mining block?");
            if (!userConfirmation) return;
        }

        const response1 = await ProofAiService.logout();
        if (response1.error) {
            showAlert(response1.error, "danger");
            return;
        }

        sessionStorage.clear();
        ipcRenderer.send("restart-app");
    };

    const handleSetServiceMachineAddr = async () => {
        setServiceMachineAddr(serviceMachineAddr.trim());
        setServiceButtondisable(true);
        hideAlert();
        const pingResponse = await ProofAiService.pingServiceMachineAddr(serviceMachineAddr);
        if (pingResponse.error) {
            showAlert(pingResponse.error, "error");
            setServiceButtondisable(false);
            return;
        }

        const response = await ProofAiService.setServiceMachineAddr(serviceMachineAddr);
        setIsServiceAddressSet(true);
        setShowModal(false);
        setServiceButtondisable(false);
    };


    useEffect(() => {
        setShowModal(true);
        //setIsServiceAddressSet(true)
    }, []);

    const routes = [
        { path: '/', element: <Login /> },
        { path: '/Login', element: <Login /> },
        { path: '/Home', element: <Home /> },
        { path: '/NewTransaction', element: <NewTransaction /> },
        { path: '/MinedTransactions', element: <MinedTransactions /> },
        { path: '/MinedTransaction', element: <MinedTransaction /> },
        { path: '/CurrentlyBlock', element: <CurrentlyBlock /> },
        { path: '/UserInstructions', element: <UserInstructions /> }
    ];


    return {
        navigate,
        ProofAiService,
        showAlert,
        hideAlert,
        isServiceAddressSet,
        setIsServiceAddressSet,
        serviceMachineAddr,
        setServiceMachineAddr,
        showModal,
        setShowModal,
        serviceButtondisable,
        setServiceButtondisable,
        handleLogout,
        handleSetServiceMachineAddr,
        routes
    }

}

export default useApp;


