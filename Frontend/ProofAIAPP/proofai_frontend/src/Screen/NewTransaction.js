import React from 'react'
import IPFS from '../Components/Ipfs'
import MiningTransaction from '../Components/MiningTransaction'
import { useProofAiService } from '../ProofaiServiceContext'
import { useAlert } from '../Context/AlertContext'
import { ClipLoader } from 'react-spinners'
const NewTransaction = () => {

    const [ipfsCont, setipfsCont] = React.useState(false)
    const [mineTransaction, setmineTransaction] = React.useState(false)
    const [modelCID, setmodelCID] = React.useState('')
    const [datasetCID, setDatasetCID] = React.useState('')
    const [puzzleStrength, setPuzzleStrength] = React.useState(0)
    const [responseWait, setResponseWait] = React.useState(false)


    const proofAiService = useProofAiService()
    const { showAlert, hideAlert } = useAlert();
    const [transactionJson, settransactionJson] = React.useState({})


    const handleIpfs = () => {
        hideAlert();
        setipfsCont(!ipfsCont)
    }
    const handleMine = async () => {
        if (modelCID === '' || datasetCID === '') {
            showAlert('Please Enter Model and Dataset CID', 'error')
            return
        }
        hideAlert()

        const Roleresponse = await proofAiService.getRole()
        if (Roleresponse.error) {
            alert(Roleresponse.error, 'error')
            return
        }
        if (Roleresponse.role !== 'Miner') {
            alert('Only miners can process transactions. Please switch your role to Miner to proceed.', 'error')
            return
        }



        setResponseWait(true)
        const response = await proofAiService.newTransaction({ modelCID, datasetCID, puzzleStrength })
        setResponseWait(false)
        if (response.error) {
            alert.show(response.error, 'error')
            return
        }

        settransactionJson(response.transaction)
        setmineTransaction(!mineTransaction)
        setDatasetCID('')
        setmodelCID('')
    }

    const handleBackButton = () => {
        hideAlert();
        window.history.back()
    }


    const handlePuzzleStrengthChange = (event) => {
        setPuzzleStrength(event.target.value);
    };

    return (
        <div style={styles.container} >


            {ipfsCont && <IPFS setipfsCont={setipfsCont} />}
            {mineTransaction && <MiningTransaction setmineTransaction={setmineTransaction} transactionJson={transactionJson} />}
            {responseWait && <ClipLoader color={'#123abc'} loading={true} size={50} />}

            {(!ipfsCont && !mineTransaction) && <div style={styles.box} >
                <h1 style={styles.label} >New Transaction</h1>
                <label> Enter CID of Model  <input style={styles.input} onChange={(event) => { setmodelCID(event.target.value) }} ></input> </label>
                <label> Enter CID of Dataset  <input style={styles.input} onChange={(event) => { setDatasetCID(event.target.value) }} ></input> </label>


                <div style={styles.buttonContianer} >
                    <button style={styles.button} onClick={handleMine} >Mine Transaction</button>
                    <button style={styles.button} onClick={handleIpfs}  >Need to Upload First on IPFS</button>
                </div>
                <div style={styles.buttonContianer} >
                    <button style={styles.button} onClick={handleIpfs}  >Transaction Instruction</button>

                    <button style={{ ...styles.button }} onClick={handleBackButton}   >back</button>
                </div>
            </div>}

        </div>
    )

}

export default NewTransaction

const styles = {
    slider: {
        width: "50%",
        height: "8px",
        borderRadius: "5px",
        background: "#ddd",
        outline: "none",
        appearance: "none",
    },

    container: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: "100%",
        height: "100vh",
        backgroundColor: "#f4f6f7",
        flexDirection: 'column',
    },
    label: {
        fontSize: '27px',
        fontWeight: 'bold',
        color: '#333',
        textAlign: 'center',
    },

    box: {
        display: 'flex',
        flexDirection: 'column',
        gap: '10px',
        padding: '20px',
        backgroundColor: "rgb(152, 156, 158)",
        borderRadius: '10px',
        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
        width: '50%',
        border: '1px solid #bc8c49',
    },

    input: {
        display: 'block',
        width: '100%',
        padding: '8px',
        marginTop: '5px',
        borderRadius: '5px',
        border: '1px solid #ccc',
        boxSizing: 'border-box',
    },

    buttonContianer: {
        gap: '10px',
        display: 'flex',
        justifyContent: 'space-around',
    },
    button: {
        padding: '10px 20px',
        border: 'none',
        borderRadius: '5px',
        backgroundColor: 'rgb(15, 67, 89)',
        color: '#fff',
        cursor: 'pointer',
        fontSize: '14px',
        width: '50%',
        transition: 'background-color 0.3s ease',
        border: '2px solid #bc8c49',
    },
    buttonHover: {
        backgroundColor: '#0056b3',
    },

}