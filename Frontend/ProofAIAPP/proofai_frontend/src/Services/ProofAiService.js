import axios from 'axios';
import path from 'path-browserify';

class ProofAiService {


    constructor() {
        this.baseUrl = 'http://localhost:8080/api';
        this.isUserLoggedIn = false;
    }

    async getPublicKey() {
        try {
            const response = await axios.get(`${this.baseUrl}/Pubkey`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data;
        } catch (error) {
            return error.response.data;
        }
    }

    async setTransactionsStorageLocation() {
        try {
            // File name and simulated path (for demonstration)
            const fileName = 'Transaction.json';
            const filePath = path.join('/simulated/directory', fileName);

            // Initialize file data
            const fileData = JSON.stringify([]);

            // Save the file data in localStorage or for download
            localStorage.setItem(fileName, fileData);

            alert(`File created and stored virtually at: ${filePath}`);
            return filePath;
        } catch (error) {
            console.error('Error setting up transactions storage:', error);
            throw error;
        }
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

    async getServiceMachineAddr() {
        try {
            const response = await axios.get(`${this.baseUrl}/GetServiceMachineIP`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data.serviceMachineIP;
        } catch (error) {
            return error.response.data;
        }
    }

    async pingServiceMachineAddr(serviceMachineAddr) {
        try {
            const response = await axios.get(`${"http://" + serviceMachineAddr}/fetch`);
            alert(response);
            return response;
        } catch (error) {

            if (error.response) {
                return error.response.data;
            }
            else {
                return { error: "Service Machine is not reachable" };
            }
        }

    };

    async logout() {
        if (this.isUserLogoin === false) {
            return { error: "Please Login First" };
        }

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

        if (this.isUserLogoin === false) {
            return { error: "Please Login First" };
        }

        try {
            const response = await axios.get(`${this.baseUrl}/getCurrentlyMinBlock`, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            return response.data;
        } catch (error) {
            if (error.response) {
                return error.response.data;
            } else {
                return { error: "Please Login First" };
            }
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
    
            if (response.data.login === "Success") {
                this.isUserLoggedIn = true;  // Fixed typo
            }
    
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
        if (!this.isUserLoggedIn) {  
            return { error: "Please Login First" };
        }
    
        const response = await this.getCurrentlyMinBlock();

        if (response.block !== "null" && response.block !== null) {  // Fixed condition
                return { error: "Cannot change role to Validator while mining is in progress" };
        }
        

        
    
        try {
            const params = new URLSearchParams();
            params.append('role', role);
            
            const response = await axios.post(`${this.baseUrl}/setRole`, params, {
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
            });
    
            return response.data;
        } catch (error) {
            console.error("API call failed:", error);
    
            return { 
                error: error.response?.data || error.message || "Unknown Error"
            };
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
        const machineAddr = await this.getServiceMachineAddr();
        const url = `http://${machineAddr}/upload`;

        try {
            const response = await axios.post(url, data, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            // Assuming the response has a "cid" field, adjust according to your response
            return response.data.message;
        } catch (error) {
            // Log the error if something goes wrong
            console.error("Upload failed:", error);
            alert(error)
            return error.response ? error.response.data : { error: "An unknown error occurred" };
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