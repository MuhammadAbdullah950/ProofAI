package main

/*
In this file we have the functions used in the application.
1. writeTransaction: function to write a transaction to the buffer
2. readTransaction: function to read a transaction from the buffer
3. broadcastTransaction: function to broadcast a transaction to all miners
4. signTransaction: function to sign a transaction
5. downloadFromIPFS: function to download files from IPFS
6. transactionHash: function to compute the hash of a transaction
7. verifyTransaction: function to verify the signature of a transaction
8. hashStruct: function to compute the hash of a struct
9. findBlockBy_Nonce_From: function to find a block by nonce and from address
10. IncomingBlockVerfication: function to verify an incoming block
11. generateBlock: function to generate a block
12. generateSalt: function to generate a salt
13. PoW: function to perform Proof of Work
14. RemoveByIndex: function to remove transactions by index
15. BlockMining: function to mine a block
16. MineTransaction: function to mine a transaction
17. InsertBlockInledgerFile: function to insert a block in the ledger file

*/

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"
)

/*
writeTransaction is a function to write a transaction to the buffer
 1. miner: miner object
 2. transaction: transaction object
 3. return: error
    Json format of the transaction is written to the buffer
*/
func writeTransaction(miner *Miner, transaction interface{}) error {

	if miner.write == nil {
		return fmt.Errorf(("Error: write buffer is nil"))
	}

	jsonData, err := json.Marshal(transaction)

	if err != nil {
		fmt.Printf("Error converting transaction to json format: %v", err)
	}
	_, err = miner.write.WriteString(string(jsonData) + "\n")
	if err != nil {
		return fmt.Errorf("Error writing transaction to buffer : %v", err)
	}
	miner.write.Flush()

	return nil
}

/*
readTransaction is a function to read a transaction from the buffer
 1. miner: miner object
    Read the transaction from the buffer and convert it to a transaction object
    If the transaction is a block, verify the transactions in the block and add it to the ledger if valid and broadcast it to all miners
*/
func readTransaction(miner *Miner) {

	for {
		jsonStr, err := miner.read.ReadString('\n')

		if err != nil {

			if err.Error() == "EOF" {
				log.Printf("Connection closed by miner: %v", err)
				return
			} else {
				log.Printf("Error from reading transaction %v : ", err)
				return
			}
		} else {

			var data map[string]interface{}
			err = json.Unmarshal([]byte(jsonStr), &data)
			if err != nil {
				log.Printf("Error converting Transaction in json format to in Transacton struct ")
			} else {

				switch data["type"] {

				case "transaction":

					var transaction Transaction
					json.Unmarshal([]byte(jsonStr), &transaction)
					if _, exists := ProofAI.receivedTransaction[transaction.Signature]; !exists {
						ProofAI.receivedTransaction[transaction.Signature] = true
						ProofAI.memPool.transactions = append(ProofAI.memPool.transactions, transaction)
						fmt.Println("Transaction Received")
						broadcastTransaction(ProofAI.Miners, transaction)
					}

				case "block":
					var block Block
					json.Unmarshal([]byte(jsonStr), &block)
					fmt.Println(time.Now())
					if _, exists := ProofAI.receivedBlock[block.TransactionsHash]; !exists {
						ProofAI.receivedBlock[block.TransactionsHash] = true
						fmt.Println("Block Received to insert in ledger")
						ProofAI.selfMiningDetail.cancel() // stop mining

						for !ProofAI.selfMiningDetail.interuptStatus {
							time.Sleep(1 * time.Second)
						}

						broadcastTransaction(ProofAI.Miners, block)
						IncomingBlockVerfication(&block)
					}
				}
			}
		}
	}
}

/*
broadcastTransaction is a function to broadcast a transaction to all miners
 1. miners: list of miners
 2. transaction: transaction object
    Write the transaction to all miners
*/
func broadcastTransaction(miners []Miner, transaction interface{}) {

	for _, miner := range miners {
		err := writeTransaction(&miner, transaction)
		if err != nil {
			fmt.Printf("error writing trnasaction to miner %s : %v", miner.conn.RemoteAddr(), err)
		}
	}
}

/*
signTransaction is a function to sign a transaction
 1. privateKey: private key of the miner
 2. transHash: hash of the transaction
    Decode the transaction hash from hex
    Ensure the hash is no longer than 32 bytes
    Sign the hash
    Encode the signature in DER format
*/
func signTransaction(privateKey *ecdsa.PrivateKey, transHash string) (string, error) {
	hashBytes, err := hex.DecodeString(transHash)
	if err != nil {
		return "", fmt.Errorf("error decoding hash: %v", err)
	}

	curveOrderSize := privateKey.PublicKey.Curve.Params().N.BitLen() / 8
	if len(hashBytes) > curveOrderSize {
		// Truncate the hash to match the curve's order size
		hashBytes = hashBytes[:curveOrderSize]
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashBytes)
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{R: r, S: s})
	if err != nil {
		return "", fmt.Errorf("error encoding signature: %v", err)
	}
	return hex.EncodeToString(signature), nil
}

/*
FileInfo is a struct to store file information
 1. Name: name of the file
 2. Hash: hash of the file
 3. Content: content of the file
*/
type FileInfo struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

/*
Response is a struct to store the response from the server
 1. Success: status of the response
 2. Message: message from the server
 3. Files: list of files
*/
type Response struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Files   []FileInfo `json:"files"`
}

/*
downloadFromIPFS is a function to download files from IPFS
 1. cid: content identifier of the file
 2. outputDir: output directory to save the files
    Create the output directory
    Prepare the request
    Check the status code
    Parse the response
    Save each file
*/
func downloadFromIPFS(cid string, outputDir string) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	serviceMachineURl := "http://" + ProofAI.selfMiningDetail.serviceMachineAddr

	data := bytes.NewBufferString(fmt.Sprintf("cid=%s", cid))
	resp, err := http.Post(serviceMachineURl+"/fetch", "application/x-www-form-urlencoded", data)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	// Parse response
	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Message)
	}

	for _, file := range response.Files {
		filePath := filepath.Join(outputDir, file.Name)

		// Write the content directly without base64 decoding
		err = os.WriteFile(filePath, []byte(file.Content), 0644)
		if err != nil {
			return fmt.Errorf("failed to save file %s: %v", file.Name, err)
		}
	}

	return nil
}

/*
transactionHash is a function to compute the hash of a transaction
 1. transaction: transaction object
    Ensure the transaction object is not nil
    Ensure the 'from' field is not empty
    Ensure the input fields are not empty
    Compute the SHA-256 hash of the transaction data
*/
func transactionHash(transaction *Transaction) string {
	if transaction == nil {
		log.Fatalf("Error: transaction object is nil")
		return ""
	}

	if transaction.From == "" {
		log.Fatalf("Error: 'from' field in transaction is empty")
		return ""
	}
	if transaction.Input_dataSet == "" || transaction.Input_model == "" {
		log.Fatalf("Error: Input fields (dataSet or model) are empty")
		return ""
	}

	txdata := fmt.Sprintf("%s%d%s%s", transaction.From, transaction.Nonce, transaction.Input_dataSet, transaction.Input_model)
	hash := sha256.Sum256([]byte(txdata))
	return hex.EncodeToString(hash[:])
}

/*
verifyTransaction is a function to verify the signature of a transaction
 1. publicKey: public key of the sender
 2. transaction: transaction object
    Compute the hash of the transaction data
    Decode the signature
    Verify the signature
*/
func verifyTransaction(publicKey *ecdsa.PublicKey, transaction *Transaction) (bool, error) {

	txData := fmt.Sprintf("%s%d%s%s", transaction.From, transaction.Nonce, transaction.Input_dataSet, transaction.Input_model)

	hash := sha256.New()
	hash.Write([]byte(txData))
	transactionHash := hash.Sum(nil)

	signatureBytes, err := hex.DecodeString(transaction.Signature)
	if err != nil {
		return false, fmt.Errorf("error decoding signature: %v", err)
	}

	var sig struct {
		R, S *big.Int
	}
	_, err = asn1.Unmarshal(signatureBytes, &sig)
	if err != nil {
		return false, fmt.Errorf("error unmarshaling signature: %v", err)
	}

	valid := ecdsa.Verify(publicKey, transactionHash, sig.R, sig.S)
	if !valid {
		return false, fmt.Errorf("signature verification failed")
	}

	return true, nil
}

/*
hashStruct is a function to compute the hash of a struct
 1. data: struct data
    Serialize the struct data to JSON format
    Compute the SHA-256 hash of the JSON data
*/
func hashStruct(data interface{}) (string, error) {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Error serializing struct %v")
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}

/*
findBlockBy_Nonce_From is a function to find a block by nonce and from address
 1. nonce: nonce of the transaction
 2. from: address of the sender
    Ensure the nonce is valid
    Ensure the from address is not empty
    Iterate through the transactions in the self-mined block
    Check if the transaction with the given nonce and from address exists
*/
func findBlockBy_Nonce_From(nonce int, from string) bool {
	if nonce < 0 || from == "" {
		fmt.Println("Invalid nonce or from address provided.")
		return false
	}

	for _, transactions := range ProofAI.selfMiningDetail.CurrentlyMineBlock.Transactions {
		if transactions.From == from && transactions.Nonce == nonce {
			return true
		}
	}
	return false
}

/*
IncomingBlockVerfication is a function to verify an incoming block
 1. block: incoming block
    Verify the transactions in the block
    Generate the transaction hash for the self-mined block
    Generate the hash for the previous block
    Compare the hashes
    Update the self-mined block properties
    Generate hashes for comparison
    Compare hashes
    Add the incoming block to the ledger if valid and update the ledger
*/
func IncomingBlockVerfication(block *Block) {
	fmt.Println("Incoming Block Verification started")

	transVerification := true

	for i, transaction := range block.Transactions {
		// Ensure index is within the range of self-mined block transactions
		if i >= len(ProofAI.selfMiningDetail.CurrentlyMineBlock.Transactions) {
			fmt.Printf("No corresponding transaction in self-mined block for index %d\n", i)
			continue
		}

		fmt.Printf("Block Transaction nonce: %d, From: %s\n", transaction.Nonce, transaction.From)

		transactionExist := findBlockBy_Nonce_From(transaction.Nonce, transaction.From)
		if !transactionExist {
			fmt.Printf("Transaction with Nonce: %d, From: %s needs mining.\n", transaction.Nonce, transaction.From)
			MineTransaction(&transaction, &ProofAI.selfMiningDetail.CurrentlyMineBlock)
			transVerification = false
		}
	}

	if transVerification {
		fmt.Println("Transaction verification completed")

		transHash, err := hashStruct(ProofAI.selfMiningDetail.CurrentlyMineBlock.Transactions)
		if err != nil {
			fmt.Printf("Error constructing block hash of transactions: %v\n", err)
			return
		}

		ProofAI.selfMiningDetail.CurrentlyMineBlock.TransactionsHash = transHash

		blockSize := len(ProofAI.ledger.blocks)
		var prevBlockHash string

		if blockSize != 0 {
			prevBlockHash, err = hashStruct(ProofAI.ledger.blocks[blockSize-1])
			if err != nil {
				fmt.Println("Error constructing previous block hash")
				return
			}

			ProofAI.selfMiningDetail.CurrentlyMineBlock.BlockNum = ProofAI.ledger.blocks[blockSize-1].BlockNum + 1
		} else {
			prevBlockHash = "000GENESISBLOCKID"
			ProofAI.selfMiningDetail.CurrentlyMineBlock.BlockNum = 1
		}

		ProofAI.selfMiningDetail.CurrentlyMineBlock.Prev_Hash = prevBlockHash
		ProofAI.selfMiningDetail.CurrentlyMineBlock.ProposerId = ProofAI.selfMiningDetail.pubKeyStr
		ProofAI.selfMiningDetail.CurrentlyMineBlock.Difficulty = ProofAI.difficultyLevel
	}

	ProofAI.selfMiningDetail.CurrentlyMineBlock.Salt = block.Salt
	ProofAI.selfMiningDetail.CurrentlyMineBlock.Prev_Hash = block.Prev_Hash
	ProofAI.selfMiningDetail.CurrentlyMineBlock.Proof = block.Proof
	ProofAI.selfMiningDetail.CurrentlyMineBlock.TimeStamp = block.TimeStamp
	ProofAI.selfMiningDetail.CurrentlyMineBlock.EventEmit = block.EventEmit
	ProofAI.selfMiningDetail.CurrentlyMineBlock.Type = block.Type
	ProofAI.selfMiningDetail.CurrentlyMineBlock.ProposerId = block.ProposerId
	ProofAI.selfMiningDetail.CurrentlyMineBlock.BlockNum = block.BlockNum
	ProofAI.selfMiningDetail.CurrentlyMineBlock.Difficulty = block.Difficulty

	currentBlockHash, err := hashStruct(ProofAI.selfMiningDetail.CurrentlyMineBlock)
	if err != nil {
		fmt.Printf("Error during hash generation: %v\n", err)
		return
	}

	incomingBlockHash, err := hashStruct(block)
	if err != nil {
		fmt.Printf("Error during hash generation for incoming block: %v\n", err)
		return
	}

	// Compare hashes
	if currentBlockHash != incomingBlockHash {
		fmt.Println("Error: Incoming block hash does not match. Incoming block is invalid.")
		return
	}

	fmt.Println("Incoming block is valid and will now be added to the ledger.")
	ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, *block)

}

/*
MineTransaction is a function to mine a transaction
 1. transaction: transaction object
 2. block: block object
    Compute the hash of the transaction
    Sign the transaction
*/
func generateBlock(trans_list []Transaction, wg *sync.WaitGroup, ctx context.Context) {

	defer wg.Done()
	fmt.Printf("New Block start to mining. It contains %d transactions.\n", len(trans_list))

	ProofAI.selfMiningDetail.CurrentlyMineBlock = Block{}
	ProofAI.CurrentlyMineBlock = &ProofAI.selfMiningDetail.CurrentlyMineBlock

	ProofAI.CurrentlyMineBlock.Transactions = trans_list

	// Set block properties
	blockSize := len(ProofAI.ledger.blocks)
	var prev_blockHash string
	var err error
	if blockSize != 0 {
		ProofAI.CurrentlyMineBlock.BlockNum = ProofAI.ledger.blocks[blockSize-1].BlockNum + 1
		// Hash previous block
		prev_blockHash, err = hashStruct(ProofAI.ledger.blocks[blockSize-1])
		if err != nil {
			fmt.Println("Error constructing previous block hash.")
			return
		}
	} else {
		prev_blockHash = "00000GENESISBLOCKID"
		ProofAI.CurrentlyMineBlock.BlockNum = 1
	}

	ProofAI.CurrentlyMineBlock.Prev_Hash = prev_blockHash
	ProofAI.CurrentlyMineBlock.ProposerId = ProofAI.selfMiningDetail.pubKeyStr
	ProofAI.CurrentlyMineBlock.Difficulty = ProofAI.difficultyLevel

	ProofAI.currentlyMiningBlockForUser = *ProofAI.CurrentlyMineBlock
	ProofAI.CurrentlyMineBlock.Transactions = nil

	for _, transaction := range trans_list { // Process transactions
		// Check if context is canceled before mining each transaction
		select {
		case <-ctx.Done():
			fmt.Printf("Mining interrupted during transaction processing for block %d.\n", ProofAI.CurrentlyMineBlock.BlockNum)
			ProofAI.selfMiningDetail.interuptStatus = true
			return
		default:
			MineTransaction(&transaction, ProofAI.CurrentlyMineBlock)
		}
	}

	trans_hash, err := hashStruct(ProofAI.CurrentlyMineBlock.Transactions)
	if err != nil {
		fmt.Printf("Error constructing block hash of transactions: %v\n", err)
		return
	}
	ProofAI.CurrentlyMineBlock.TransactionsHash = trans_hash

	fmt.Printf("Block number: %d\n", ProofAI.CurrentlyMineBlock.BlockNum)

	// Check context before starting proof of work
	if ctx.Err() != nil {
		fmt.Printf("Interrupted mining block %d before Proof of Work (PoW) started.\n", ProofAI.CurrentlyMineBlock.BlockNum)
		ProofAI.selfMiningDetail.interuptStatus = true
		return
	}

	// Perform Proof of Work
	err = PoW(ProofAI.CurrentlyMineBlock, ctx)
	if err != nil {
		fmt.Printf("Error during Proof of Work for block: %v\n", err)
		return
	}

	// Finalize block details
	ProofAI.CurrentlyMineBlock.Type = "block"
	ProofAI.CurrentlyMineBlock.TimeStamp = time.Now().Format(time.RFC3339)

	// Log block mining success
	fmt.Println("Block mined successfully. Broadcasting to all miners...")

	// Add to receivedBlock and broadcast
	ProofAI.receivedBlock[ProofAI.CurrentlyMineBlock.TransactionsHash] = true

	broadcastTransaction(ProofAI.Miners, ProofAI.CurrentlyMineBlock)
	InsertBlockInledgerFile("ledger.json", ProofAI.CurrentlyMineBlock)

	// Log broadcast completion
	fmt.Printf("Block broadcasted successfully at %s.\n", time.Now().Format(time.RFC3339))
	ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, *ProofAI.CurrentlyMineBlock)
	ProofAI.CurrentlyMineBlock = nil
	ProofAI.currentlyMiningBlockForUser = Block{}
}

/*
generateSalt is a function to generate a salt
 1. size: size of the salt
    Generate a random salt of the given size
    Encode the salt to hex format
    Return the salt and error
*/
func generateSalt(size uint) (string, error) {
	salt := make([]byte, size)

	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("Failed to generate Salt:  %v", err)
	}
	return hex.EncodeToString(salt), err
}

/*
PoW is a function to perform Proof of Work
 1. block: block object
    Set the block size and salt size
    Set the difficulty level
    Generate the required prefix
    Generate the salt
    Compute the hash of the block
    Check if the hash has the required prefix
    If the hash has the required prefix, set the salt and return nil
    If the context is canceled, set the interupt status and return an error
*/
func PoW(block *Block, ctx context.Context) error {
	fmt.Printf("Pow of work start")
	blockSize := uint(102400)
	saltSize := blockSize - uint(unsafe.Sizeof(block))
	block.Difficulty = 3
	requiredPrefix := strings.Repeat("0", block.Difficulty)
	for {
		salt, err := generateSalt(saltSize)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		block.Salt = salt
		blockHash, err := hashStruct(block)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		// fmt.Println(blockHash)
		if strings.HasPrefix(blockHash, requiredPrefix) {
			//fmt.Printf("Block : %+v", block)
			block.Salt = salt
			fmt.Printf(" BlockHash : %s\n", blockHash)
			return nil
		}

		if ctx.Err() != nil {
			ProofAI.selfMiningDetail.interuptStatus = true
			return fmt.Errorf("Interup during POW")
		}
	}
}

/*
RemoveByIndex is a function to remove transactions by index
 1. transaction: list of transactions
 2. indices: list of indices to remove
    Iterate through the indices in reverse order
    Remove the transaction at the given index
*/
func RemoveByIndex(transaction []Transaction, indices []int) []Transaction {
	for i := len(indices) - 1; i >= 0; i-- {
		index := indices[i]
		if index < 0 || index >= len(transaction) {
			continue // Skip invalid indices
		}
		transaction = append(transaction[:index], transaction[index+1:]...)
	}
	return transaction
}

/*
BlockMining is a function to mine a block
 1. ctx: context object
    Check if the ledger has blocks
    Get the last block
    Parse the timestamp of the last block
    Get the current time and calculate the difference
    If the difference is greater than 5 minutes, start mining
    Get only 2 transactions from the mempool if available, otherwise start mining with one transaction
*/
func BlockMining(ctx context.Context) {
	fmt.Println("MemPool Start To Receive Transactions")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Block mining stopped")
			return
		default:

			if len(ProofAI.ledger.blocks) > 0 {
				lastBlock := ProofAI.ledger.blocks[len(ProofAI.ledger.blocks)-1]
				lastBlockTime, err := time.Parse(time.RFC3339, lastBlock.TimeStamp)
				if err != nil {
					fmt.Printf("Error parsing time: %v\n", err)
					return
				}
				// get current time and calculate difference and if dfference is greater than 5 minutes then start mining
				currentTime := time.Now()
				diff := currentTime.Sub(lastBlockTime)
				if diff.Minutes() > 2 {
					// get only 2 transactions from mempool if available otherwise start mining with one transaction
					if len(ProofAI.memPool.transactions) >= 2 {

						var wg sync.WaitGroup
						ProofAI.selfMiningDetail.context, ProofAI.selfMiningDetail.cancel = context.WithCancel(context.Background())
						wg.Add(1)

						go generateBlock(ProofAI.memPool.transactions[:2], &wg, ProofAI.selfMiningDetail.context)
						wg.Wait()
						removeList := []int{0, 1}
						ProofAI.memPool.transactions = RemoveByIndex(ProofAI.memPool.transactions, removeList)
					} else if len(ProofAI.memPool.transactions) == 1 {
						var wg sync.WaitGroup
						ProofAI.selfMiningDetail.context, ProofAI.selfMiningDetail.cancel = context.WithCancel(context.Background())
						wg.Add(1)
						go generateBlock(ProofAI.memPool.transactions[:1], &wg, ProofAI.selfMiningDetail.context)
						wg.Wait()
						removeList := []int{0}
						ProofAI.memPool.transactions = RemoveByIndex(ProofAI.memPool.transactions, removeList)
					}
				} else {
					time.Sleep(1 * time.Minute)
				}
			} else {
				if len(ProofAI.memPool.transactions) >= 2 {

					var wg sync.WaitGroup
					ProofAI.selfMiningDetail.context, ProofAI.selfMiningDetail.cancel = context.WithCancel(context.Background())
					wg.Add(1)

					go generateBlock(ProofAI.memPool.transactions[:2], &wg, ProofAI.selfMiningDetail.context)
					wg.Wait()
					removeList := []int{0, 1}
					ProofAI.memPool.transactions = RemoveByIndex(ProofAI.memPool.transactions, removeList)
				} else if len(ProofAI.memPool.transactions) == 1 {
					var wg sync.WaitGroup
					ProofAI.selfMiningDetail.context, ProofAI.selfMiningDetail.cancel = context.WithCancel(context.Background())
					wg.Add(1)
					go generateBlock(ProofAI.memPool.transactions[:1], &wg, ProofAI.selfMiningDetail.context)
					wg.Wait()
					removeList := []int{0}
					ProofAI.memPool.transactions = RemoveByIndex(ProofAI.memPool.transactions, removeList)
				}
			}
		}
	}
}

/*
MineTransaction is a function to mine a transaction
 1. transaction: transaction object
 2. block: block object
    Download the dataset and model from IPFS
    Execute the model
    Add the transaction to the block
*/
func MineTransaction(transaction *Transaction, block *Block) {

	fmt.Printf("Mining Transaction : %+v\n", transaction.From)

	pubkey, err := hexToPublicKey(transaction.From)
	if err != nil {
		log.Fatalf("Error getting public key from hex string %v", err)
		return
	}

	signatureValidation, err := verifyTransaction(pubkey, transaction)
	if err != nil {
		log.Fatalf("Error verifying transaction %v", err)
		return
	}
	if !signatureValidation {
		log.Fatalf("Error: Transaction signature is invalid")
		return
	} else {
		fmt.Println("Transaction signature is valid")
	}

	// Determine the directory path
	currentDir := os.TempDir()
	// add time stamp with the directory
	dirPath := filepath.Join(currentDir, ProofAI.modelExecutionDir + time.Now().Format("20060102150405"))  
	fmt.Println("Directory path: ", dirPath)

	// Check if the directory exists and remove it
	if _, err := os.Stat(dirPath); err == nil {
		// Directory exists; remove it
		if err := os.RemoveAll(dirPath); err != nil {
			fmt.Printf("Error removing existing directory: %v\n", err)
			return
		}
		fmt.Println("Existing directory removed")
	} else if !os.IsNotExist(err) {
		// Unexpected error
		fmt.Printf("Error checking directory existence: %v\n", err)
		
		return
	}

	// Create a fresh directory
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		cleanDir(dirPath )
		return
	}
	fmt.Println("Directory created successfully")

	// Download dataset and model to the specified directory
	err = downloadFromIPFS(transaction.Input_dataSet, filepath.Join(dirPath, "dataset"))
	if err != nil {
		log.Fatalf("Error downloading dataset from IPFS: %v", err)
		cleanDir(dirPath )
		return
	}

	err = downloadFromIPFS(transaction.Input_model, filepath.Join(dirPath, "model"))
	if err != nil {
		log.Fatalf("Error downloading model from IPFS: %v", err)
		cleanDir(dirPath )
		return
	}

	// Execute the model and update the transaction
	modelOutput, err := modelExecution(dirPath)
	if err != nil {
		fmt.Printf("Error during model execution: %v\n", err)
		cleanDir(dirPath )
		return
	}
	transaction.Model_output = modelOutput
	block.Transactions = append(block.Transactions, *transaction)
	fmt.Println("Mining Transaction Completed")

	cleanDir(dirPath ) 
	
}

func cleanDir(dirpath string) {
	if err := os.RemoveAll(dirpath); err != nil {
		fmt.Printf("Error removing directory at cleanup: %v\n", err)
		return
	}
	fmt.Println("Temporary directory cleaned up")}


/*
userTransaction is a function to create a transaction for the user
 1. model_cid: content identifier of the model
 2. dataset_cid: content identifier of the dataset
    Create a transaction object
    Sign the transaction
    Add the transaction to the mempool
    Broadcast the transaction to all miners
    Return the transaction object
*/
func userTransaction(model_cid string, dataset_cid string) (Transaction, error) {

	transaction_ := Transaction{
		From:          ProofAI.selfMiningDetail.pubKeyStr,
		Nonce:         ProofAI.selfMiningDetail.nonce,
		Input_dataSet: dataset_cid,
		Input_model:   model_cid,
		Type:          "transaction",
	}

	transHash := transactionHash(&transaction_)
	var err error
	transaction_.Signature, err = signTransaction(ProofAI.selfMiningDetail.prvKey, transHash)
	if err != nil {
		log.Printf("Error Signing transaction: %v\n", err)
		return Transaction{}, err
	}
	ProofAI.selfMiningDetail.nonce += 1
	ProofAI.receivedTransaction[transaction_.Signature] = true
	ProofAI.memPool.transactions = append(ProofAI.memPool.transactions, transaction_)
	fmt.Println("Transaction Recieved.")
	broadcastTransaction(ProofAI.Miners, &transaction_)

	return transaction_, nil
}
