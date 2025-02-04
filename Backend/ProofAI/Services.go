package main

/*
	In this we create two  Server to listen for incoming requests from the external world and run background services to communicate with frontend and backend services.
	1-	createServerAndListenExternelWorld creates a new ProofAI object and starts the server to listen for incoming requests from the external world.
	2-	createServerAndListen creates a new ProofAI object and starts the server to listen for incoming requests from the frontend.
	3-	handleGetLatestBlock gets the latest block.
	4-	handlegetServiceMachineIP gets the IP address of the service machine.
	5-	handlegetPubkey gets the public key of the miner.
	6-	handleServiceMachineIP sets the IP address of the service machine.
	7-	handleLogout logs out the miner.
	8-	handleSetRole sets the role of the miner.
	9-	handleGetRole gets the role of the miner.
	10-	handleGetCurrentlyMiningBlock gets the currently mining block.
	11-	handleGetMinedBlocks gets the mined blocks.
	12-	handleTransactionConfirmation checks if the transaction is confirmed.
	13-	handleGenerateKey generates the public and private keys for the miner.
	14-	handleNewTransaction creates a new transaction.
	15-	handleLoginVerification verifies the login of the miner.
	16-	successfulllogin is called when the login is successful.
	17-	sendServiceLogout sends a logout request to the service machine where the miner is connected to.
*/

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
)

/*
createServerAndListenExternelWorld creates a new ProofAI object and starts the server to listen for incoming requests
REST API is used to communicate with the service machine
*/
func createServerAndListenExternelWorld() {
	fmt.Println("Server Starting for External World...")

	http.HandleFunc("/api/latestBlock", handleGetLatestBlock)     // get latest block
	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"})) // allow all origins

	ip, err := getRadminIPv4()
	fmt.Println("IP : ", ip)
	address := ip + ":8079" // bind to all interfaces on port 8079
	fmt.Printf("Listening on %s\n", address)

	err = http.ListenAndServe(address, cors(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

/*
handleGetLatestBlock gets the latest block
Output parameter : response
*/
func handleGetLatestBlock(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("Get Latest Block APi is called")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(ProofAI.ledger.blocks) == 0 {
		response := map[string]interface{}{"block": "null"}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := map[string]interface{}{"block": ProofAI.ledger.blocks[len(ProofAI.ledger.blocks)-1]}
	json.NewEncoder(w).Encode(response)

}

/*
  - createServerAndListen creates a new ProofAI object and starts the server to listen for incoming requests
    REST API is used to communicate with the service machine
*/
func createServerAndListen() {

	http.HandleFunc("/api/login", handleLoginVerification)                         // login verification
	http.HandleFunc("/api/getRole", handleGetRole)                                 // get role of the miner
	http.HandleFunc("/api/GetServiceMachineIP", handlegetServiceMachineIP)         // get public key of the miner
	http.HandleFunc("/api/Pubkey", handlegetPubkey)                                // get public key of the miner
	http.HandleFunc("/api/logout", handleLogout)                                   // logout the miner
	http.HandleFunc("/api/setRole", handleSetRole)                                 // set role of the miner
	http.HandleFunc("/api/ServiceMachineIP", handleServiceMachineIP)               // set service machine IP
	http.HandleFunc("/api/generateKeys", handleGenerateKey)                        // generate keys for the miner
	http.HandleFunc("/api/newTransaction", handleNewTransaction)                   // create new transaction
	http.HandleFunc("/api/getMinedBlocks", handleGetMinedBlocks)                   // get mined blocks
	http.HandleFunc("/api/getCurrentlyMinBlock", handleGetCurrentlyMiningBlock)    // get currently mining block
	http.HandleFunc("/api/transactionConfirmation", handleTransactionConfirmation) // transaction confirmation

	cors := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))   // allow all origins
	err := http.ListenAndServe(":8080", cors(http.DefaultServeMux)) // listen on port 8080
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

/*
handlegetServiceMachineIP gets the IP address of the service machine
Output parameter : response
*/
func handlegetServiceMachineIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"serviceMachineIP": serviceMachineAdd}
	json.NewEncoder(w).Encode(response)
}

/*
  - handlegetPubkey gets the public key of the miner
    Output parameter : response
*/
func handlegetPubkey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"pubKey": ProofAI.selfMiningDetail.pubKeyStr}
	json.NewEncoder(w).Encode(response)
}

/*
  - handleServiceMachineIP sets the IP address of the service machine
    Input parameter : service machine IP address
    Output parameter : response
*/
func handleServiceMachineIP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Post method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error parsing form data: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ip := r.FormValue("ServiceMachineaddr")
	serviceMachineAdd = ip
	fmt.Println("service add : ", serviceMachineAdd)
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"ServiceMachineIP": "Set"}
	json.NewEncoder(w).Encode(response)
}

/*
  - handleLogout logs out the miner
    Output  : response (logout success)
    logic  : close the session of the miner.
    And in close session function, it will stop the mining process of the miner and clear the session details.
*/
func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Post method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("Logout API is called")
	CloseSession()
	fmt.Println("Session Closed")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"logout": "Success"}
	json.NewEncoder(w).Encode(response)
}

/*
  - handleSetRole sets the role of the miner
    Input parameter : role
    Output parameter : response
*/
func handleSetRole(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Post method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error parsing form data: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	role := r.FormValue("role")
	ProofAI.selfMiningDetail.role = role
	fmt.Println("Role : ", ProofAI.selfMiningDetail.role)
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"role": "Set"}
	json.NewEncoder(w).Encode(response)

}

/*
  - handleGetRole gets the role of the miner
    Output parameter : response
*/
func handleGetRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"role": ProofAI.selfMiningDetail.role}
	json.NewEncoder(w).Encode(response)
	fmt.Println("Role : ", ProofAI.selfMiningDetail.role)
}

/*
  - handleGetCurrentlyMiningBlock gets the currently mining block
    Output parameter : response
    logic : If the miner is currently mining a block, then return the block details.
    Else return null.
*/
func handleGetCurrentlyMiningBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//fmt.Println("Currently Mining Block")

	if ProofAI.currentlyMiningBlockForUser.Transactions != nil {
		response := map[string]interface{}{"block": ProofAI.currentlyMiningBlockForUser}
		//
		//	fmt.Println("Response ", response)
		json.NewEncoder(w).Encode(response)
	} else {
		response := map[string]interface{}{"block": "null"}
		//	fmt.Println("Response ", response)
		json.NewEncoder(w).Encode(response)
	}
}

/*
  - handleGetMinedBlocks gets the mined blocks
    Output parameter : response
    logic : If the miner has mined blocks, then return the block details.
    Else return null.
*/
func handleGetMinedBlocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// get params value
	fiterValue := r.URL.Query().Get("filter")
	//	fmt.Println("Filter value : ", fiterValue)
	var response map[string]interface{}
	if fiterValue == "Own Transactions" {

		var blocks []Block
		for _, block := range ProofAI.ledger.blocks {
			var transactionList []Transaction
			for _, transaction := range block.Transactions {
				if transaction.From == ProofAI.selfMiningDetail.pubKeyStr {
					transactionList = append(transactionList, transaction)
				}
			}
			if len(transactionList) > 0 {
				temp_Block := block
				temp_Block.Transactions = transactionList
				blocks = append(blocks, temp_Block)
			}
		}
		if len(blocks) == 0 {
			response = map[string]interface{}{"blocks": "null"}
		} else {

			response = map[string]interface{}{"blocks": blocks}
		}

	} else {

		if len(ProofAI.ledger.blocks) == 0 {
			response = map[string]interface{}{"blocks": "null"}
		} else {

			response = map[string]interface{}{"blocks": ProofAI.ledger.blocks}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

/*
  - handleTransactionConfirmation checks if the transaction is confirmed
    Input parameters : from address, nonce
    Output parameter : response
    logic : Check if the transaction is confirmed or not.
    If the transaction is confirmed, then return "Confirmed".
    Else return "Pending".
*/
func handleTransactionConfirmation(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Get method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	From := r.URL.Query().Get("from")
	nonce := r.URL.Query().Get("nonce")

	//fmt.Println(From, nonce)
	// check if the transaction is confirmed
	for _, block := range ProofAI.ledger.blocks {
		for _, transaction := range block.Transactions {
			if transaction.From == From && strconv.Itoa(transaction.Nonce) == nonce {
				w.WriteHeader(http.StatusOK)
				response := map[string]string{"transaction": "Confirmed"}
				json.NewEncoder(w).Encode(response)
				fmt.Println("Transaction Confirmed")
				fmt.Println(From, nonce)
				//	fmt.Println("transaction")
				fmt.Println(transaction.From, transaction.Nonce)
				return
			}
		}
	}
	w.WriteHeader(http.StatusAccepted)
	response := map[string]string{"transaction": "pending"}
	json.NewEncoder(w).Encode(response)
	return
}

/*
  - handleGenerateKey generates the public and private keys for the miner
    Output parameter : response
*/
func handleGenerateKey(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Post method"}
		json.NewEncoder(w).Encode(response)
		return
	}
	var err error

	pubKey, prvKey, err := generateKeys()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error generating keys" + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	pubKeyStr, prvKeyStr := keyToHex(pubKey, prvKey)
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"pubKey": pubKeyStr, "prvKey": prvKeyStr}
	json.NewEncoder(w).Encode(response)
}

/*
  - handleNewTransaction creates a new transaction
    Input parameters : model CID, dataset CID
    output parameters : response
*/
func handleNewTransaction(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Transaction Recieved.")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Post method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error parsing form data: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	modelCID := r.FormValue("modelCID")
	datasetCID := r.FormValue("datasetCID")
	transaction, err := userTransaction(modelCID, datasetCID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error creating transaction: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"transaction": transaction}
	json.NewEncoder(w).Encode(response)
}

/*
  - handleLoginVerification verifies the login of the miner
    Input parameters : public key, private key
    output parameters : response
    logic : Verify the login of the miner.
    If the login is successful, then start a new session for the miner.
    And establish a connection with the service machine.
*/
func handleLoginVerification(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{"error": "Invalid Post method"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error parsing form data: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	pubKey := r.FormValue("PubKey")
	prvKey := r.FormValue("PrvKey")

	pubKeyDecoded, err := hexToPublicKey(pubKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error decoding public key: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	prvKeyDecoded, err := hexToPrivateKey(prvKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Error decoding private key: " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	loginValid := keyVerification(pubKeyDecoded, prvKeyDecoded)

	if loginValid {
		successfulllogin(pubKey, prvKey, pubKeyDecoded, prvKeyDecoded)
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"login": "Success"}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]string{"login": "Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

/*
  - successfulllogin is called when the login is successful
    Input parameters : public key, private key, public key decoded, private key decoded
    logic : Start a new session for the miner.
    And establish a connection with the service machine.
*/
func successfulllogin(pubKey string, prvKey string, pubKeyDecoded *ecdsa.PublicKey, prvKeyDecoded *ecdsa.PrivateKey) {

	StartNewSession()
	ProofAI.selfMiningDetail.pubKey = pubKeyDecoded
	ProofAI.selfMiningDetail.prvKey = prvKeyDecoded
	ProofAI.selfMiningDetail.pubKeyStr = pubKey
	ProofAI.selfMiningDetail.prvKeyStr = prvKey
	go establishConnection(ProofAI.connectionPort) // establish connection with the service machine to get the latest block details and start mining process
}

/*
 sendServiceLogout sends a logout request to the service machine where the miner is connected to.
*/

func sendServiceLogout() {

	res, err := http.Post("http://"+serviceMachineAdd+"/logout", "application/json", nil)
	if err != nil {
		fmt.Println("Error in sending logout request to service machine", err.Error())
	}
	defer res.Body.Close()
}
