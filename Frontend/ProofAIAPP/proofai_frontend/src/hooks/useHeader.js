import React from "react";
import { useProofAiService } from "../ProofaiServiceContext";

const useHeader = () => {
    const ProofAiService = useProofAiService();
    const [role, setRole] = React.useState(true);


    const handleRoleChange = async (e) => {
        const minerRole = e.target.value;
        const response = await ProofAiService.setRole(minerRole);
        if (response.error.data) {
            alert(response.error);
            e.target.value = "Miner";
            return;
        } else if (response.error) {
            alert("Please First Login");
            e.target.value = "Miner";
            return;
        }
        setRole(minerRole);
    }

    return {
        handleRoleChange,
        role,
    }

}
export default useHeader;
