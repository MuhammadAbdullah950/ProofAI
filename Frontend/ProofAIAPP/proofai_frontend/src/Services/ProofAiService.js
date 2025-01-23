import axios from 'axios';

class ProofAiService {

    constructor() {
        this.baseUrl = 'http://localhost:8080/api';
    }

    async setServiceMachineAddr(ServiceMachineaddr) {
        try {
            const params = new URLSearchParams();
            params.append('ServiceMachineaddr', ServiceMachineaddr);

            const response = await axios.post(`${this.baseUrl}/ServiceMachineIP`, params, {
                headers: { "Content-Type": "application/x-www-form-urlencoded" },

            });
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async pingServiceMachineAddr(serviceMachineAddr) {
        try {

            alert("http://" + serviceMachineAddr + "/machines");
            const response = await axios.get(`${"http://" + serviceMachineAddr}/fetch`);
            return response.data;
        } catch (error) {
            return error.response.data;
        }

    };


    async logout() {
        try {
            const response = await axios.post(`${this.baseUrl}/logout`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async getRole() {
        try {
            const response = await axios.get(`${this.baseUrl}/getRole`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async generateKeys() {
        try {
            const response = await axios.post(`${this.baseUrl}/generateKeys`);
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async getCurrentlyMinBlock() {
        try {
            const response = await axios.get(`${this.baseUrl}/getCurrentlyMinBlock`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async login_using_key(PubKey, PrvKey) {
        try {
            const params = new URLSearchParams();
            params.append('PubKey', PubKey);
            params.append('PrvKey', PrvKey);

            const response = await axios.post(`${this.baseUrl}/login`, params, {
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
            });

            return response.data;
        } catch (error) {
            if (error.response) {
                return error.response.data;  // Return response data from the server
            } else {
                return { error: error.message };  // Return error message if response is not available
            }
        }
    }

    async setRole(role) {

        if (role === 'Validator') {
            const response = await this.getCurrentlyMinBlock()
            if (response.block !== "null") {
                return { error: "Cannot change role to Validator while mining is in progress" }
            }
        }

        try {
            const params = new URLSearchParams();
            params.append('role', role);
            const response = await axios.post(`${this.baseUrl}/setRole`, params, {
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
            });
            return response.data;
        } catch (error) {
            return { error: error.response ? error.response.data : "Network Error" };
        }
    }

    async transactionConfirmation(from, nonce) {
        try {
            const response = await axios.get(`${this.baseUrl}/transactionConfirmation`, {
                params: {
                    from: from,
                    nonce: nonce,
                },
            });
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async uploadDataOnIPfs(data) {
        const jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiIwYjQ1NTE0OS1iYTgyLTQ2YWYtYjU3MC1mYTBmMDFkMTAwOWUiLCJlbWFpbCI6ImFiZHVsbGFoLmJoYXR0aS45OTAwMTFAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInBpbl9wb2xpY3kiOnsicmVnaW9ucyI6W3siZGVzaXJlZFJlcGxpY2F0aW9uQ291bnQiOjEsImlkIjoiRlJBMSJ9LHsiZGVzaXJlZFJlcGxpY2F0aW9uQ291bnQiOjEsImlkIjoiTllDMSJ9XSwidmVyc2lvbiI6MX0sIm1mYV9lbmFibGVkIjpmYWxzZSwic3RhdHVzIjoiQUNUSVZFIn0sImF1dGhlbnRpY2F0aW9uVHlwZSI6InNjb3BlZEtleSIsInNjb3BlZEtleUtleSI6ImU0NmJmYTI0NGRiMmUwNTM0NWM1Iiwic2NvcGVkS2V5U2VjcmV0IjoiZmM2NTFhM2NlNjE3NDIwYzgzNzU5OWRkYTUyNDFlOWM3MzI4ODA5Y2ExMzVlOTZlODJlMWViYTEyMDY3ZjNjMSIsImV4cCI6MTc2Njg1NTI0OH0.kcOg4gQYOnMataWvBDbRSYHZITADTU8AYaNVz3Yjyf8"
        // process.env.JWT_TOKEN;
        const url = "https://api.pinata.cloud/pinning/pinFileToIPFS"
        //process.env.IPFS_NODE_URL;
        try {
            const response = await axios.post(url, data, {
                headers: {
                    'Authorization': `Bearer ${jwt}`,
                },
            });

            return response.data.IpfsHash;

        } catch (error) {
            return error.response.data
        }
    }

    async getFolderFromIPFS(cid) {
        try {
            const response = await fetch(`http://20.28.18.86:8081/fetch?cid=${cid}`);

            if (!response.ok) {
                throw new Error(`Error fetching folder: ${response.statusText}`);
            }

            const folderBlob = await response.blob();

            const url = window.URL.createObjectURL(folderBlob);

            const link = document.createElement("a");
            link.href = url;
            link.download = `${cid}.zip`;
            link.click();

            window.URL.revokeObjectURL(url);

            console.log("Folder downloaded successfully");
        } catch (error) {
            console.error("Failed to fetch folder from IPFS:", error);
        }
    }

    async newTransaction(data) {
        try {
            const params = new URLSearchParams();
            params.append("modelCID", data.modelCID);
            params.append("datasetCID", data.datasetCID);

            const response = await axios.post(`${this.baseUrl}/newTransaction`, params, {
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
            });

            return response.data; // Return only the data from the response.
        } catch (error) {
            console.error("Error in newTransaction:", error); // Log the error for debugging.

            // Return a fallback error response or a default value if error.response is undefined.
            return error.response?.data || { error: "An unknown error occurred" };
        }
    }

    async transactionInstructions() {
    }

    async getTransactionStatus(id) {
    }

    async getMinedBlocks(filterValue) {
        try {
            const response = await axios.get(`${this.baseUrl}/getMinedBlocks`,
                {
                    params: {
                        filterValue: filterValue,
                    },
                    headers: {
                        "Content-Type": "application/json",
                    },
                });
            return response.data;

        } catch (error) {
            return { error: error.response ? error.response.data : "Network Error" };
        }

    }

    async getCurrentlyMinTransaction() {
        try {
            const response = await axios.get(`${this.baseUrl}/getCurrentlyMinTransaction`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data;

        } catch (error) {
            return { error: error.response ? error.response.data : "Network Error" };
        }
    }

    async setMinerRole(data) {
    }

}

export default ProofAiService;