import React from 'react'
import IPFS from '../Components/Ipfs'
import MiningTransaction from '../Components/MiningTransaction'
import { useProofAiService } from '../ProofaiServiceContext'
import { useAlert } from '../Context/AlertContext'
import { ClipLoader } from 'react-spinners'
import { useNavigate } from 'react-router-dom'

const useNewTransaction = () => {

    const [ipfsCont, setipfsCont] = React.useState(false)
    const [mineTransaction, setmineTransaction] = React.useState(false)
    const [modelCID, setmodelCID] = React.useState('')
    const [datasetCID, setDatasetCID] = React.useState('')
    const [puzzleStrength, setPuzzleStrength] = React.useState(0)
    const [responseWait, setResponseWait] = React.useState(false)


    const navigate = useNavigate()
    const proofAiService = useProofAiService()
    const { showAlert, hideAlert } = useAlert();
    const [transactionJson, settransactionJson] = React.useState({})

    const handleIpfs = () => {
        hideAlert();
        setipfsCont(!ipfsCont)
    }

    const handleUserInstructions = () => {
        hideAlert();
        navigate('/UserInstructions')
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

    return {
        ipfsCont,
        setipfsCont,
        mineTransaction,
        setmineTransaction,
        setmodelCID,
        setDatasetCID,
        responseWait,
        transactionJson,
        handleMine,
        handleIpfs,
        handleUserInstructions,
        handleBackButton,

    }

}

export default useNewTransaction
