import React, { useEffect, useState } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import Login from './Screen/Login';
import Home from './Screen/Home';
import Header from './Components/Header';
import NewTransaction from './Screen/NewTransaction';
import MinedTransactions from './Screen/MinedTransactions';
import MinedTransaction from './Screen/MinedTransaction';
import CurrentlyBlock from "./Screen/CurrentlyBlock";
import { useProofAiService } from "./ProofaiServiceContext";
import { useAlert } from "./Context/AlertContext";
import path from 'path-browserify';
import UserInstructions from './Screen/UserInstructions';

const { ipcRenderer } = window.require("electron");


function App() {
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

  return (
    <>
      {showModal && (
        <div className="modal">
          <div className="modal-content">
            <h2>Enter Service Machine Address </h2>
            <p>Address Eg: 127.0.0.1:0000</p>
            <input
              type="text"
              placeholder="Enter Address"
              value={serviceMachineAddr}
              onChange={(e) => setServiceMachineAddr(e.target.value)}
            />
            <div className="modal-actions">
              <button onClick={handleSetServiceMachineAddr} disabled={isServiceAddressSet}  >Submit</button>
            </div>
          </div>
        </div>
      )}

      {isServiceAddressSet && (
        <div style={styles.container}>
          <Header handleLogout={handleLogout} />
          <div style={styles.content}>
            <Routes>
              {routes.map(({ path, element }) => (
                <Route key={path} path={path} element={element} />
              ))}
            </Routes>
          </div>
        </div>
      )}

      <style jsx>{`
        .modal {
          position: fixed;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          display: flex;
          justify-content: center;
          align-items: center;
          background: rgba(0, 0, 0, 0.5);
        }
        .modal-content {
          background: white;
          padding: 20px;
          border-radius: 8px;
          text-align: center;
        }
        .modal-actions button {
          margin: 5px;
        }
      `}</style>
    </>
  );
}

export default App;

const styles = {
  container: {
    width: "100%",
    minHeight: "100vh",
    display: "flex",
    flexDirection: "column",
    backgroundColor: "#f9f5f5",
  },
  content: {
    flex: 1,
    overflowY: "auto",
  },
};
