package main

/*
In this file We have defined the functions to handle transactions and blocks in the blockchain network.
-FileInfo is a struct to store file information
-Response is a struct to store the response from the server
-Functions:
	1. writeTransaction: function to write a transaction to the buffer
	2. readTransaction: function to read a transaction from the buffer
	3. broadcastTransaction: function to broadcast a transaction to all miners
	4. signTransaction: function to sign a transaction
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
	18. cleanDir: function to clean up a directory
	19. userTransaction: function to create a transaction for the user
*/

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	defer miner.conn.Close()
	for ProofAI.selfMiningDetail.connectionAlive {
		jsonStr, err := miner.read.ReadString('\n')

		if err != nil {

			for i, m := range ProofAI.Miners {
				if m == *miner {
					// apply mutex lock below is code
					ProofAI.selfMiningDetail.mu.Lock()
					ProofAI.Miners = append(ProofAI.Miners[:i], ProofAI.Miners[i+1:]...)
					ProofAI.selfMiningDetail.mu.Unlock()
					fmt.Println("Miner removed from the list")
					break
				}
			}
			return
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
						broadcastTransaction(&ProofAI.Miners, transaction)
						if ProofAI.selfMiningDetail.role == "Miner" {
							ProofAI.memPool.transactions = append(ProofAI.memPool.transactions, transaction)
							fmt.Println("Transaction Received For mining")
						} else {
							fmt.Println("Transaction Received But not for mining")
						}
					}

				case "block":
					var block Block
					json.Unmarshal([]byte(jsonStr), &block)
					fmt.Println(time.Now())
					if _, exists := ProofAI.receivedBlock[block.TransactionsHash]; !exists {
						ProofAI.receivedBlock[block.TransactionsHash] = true
						fmt.Println("Block Received to insert in ledger")
						if ProofAI.CurrentlyMineBlock != nil {

							if len(ProofAI.ledger.blocks) == 0 || block.BlockNum > ProofAI.ledger.blocks[len(ProofAI.ledger.blocks)-1].BlockNum {
								fmt.Println(time.Now())
								if ProofAI.selfMiningDetail.role == "Miner" {
									ProofAI.selfMiningDetail.cancel() // it will stop the mining of current block and not move to the next block unitl
									for !ProofAI.selfMiningDetail.interuptStatus {
										time.Sleep(1 * time.Second)
									}
								}
								fmt.Println("Block Received to insert in ledger 2 ")
								fmt.Println(time.Now())

								broadcastTransaction(&ProofAI.Miners, block)
								fmt.Println("Block Received to insert in ledger 3 ")

								IncomingBlockVerfication(&block)
							} else {
								fmt.Println("Block is already mined and inserted in ledger")
							}
						} else {
							fmt.Println("Block Skipped due to soon connection open")
						}
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
func broadcastTransaction(miners *[]Miner, transaction interface{}) {
	i := 0
	for i < len(*miners) {

		miner := (*miners)[i]
		err := writeTransaction(&miner, transaction)
		if err != nil {
			fmt.Printf("Error writing transaction to miner %v:\n", miner)
		}
		i++
	}
}

/*
get Latest block and now compare with ledger check hash , and check hash of all miners
block and choose having highest same hash block
*/
func updateLedger() {
	// Loop through all miners
	var latestMinersBlock []Block
	for _, miner := range ProofAI.Miners {
		IP := strings.Split(miner.conn.RemoteAddr().String(), ":")[0]
		url := "http://" + IP + ":8079/api/latestBlock"
		fmt.Println("url", url)
		res, err := http.Get(url)

		if err != nil {
			fmt.Printf("Error updating ledger from %s: %v\n", IP, err)
			continue
		}
		defer res.Body.Close()

		// Read response body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body from %s: %v\n", IP, err)
			continue
		}

		// Parse the response
		var response map[string]interface{}
		//	fmt.Println("Response from miner", string(body))

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Printf("Error unmarshalling response from %s: %v\n", IP, err)
			continue
		}

		// Extract the block from the response
		if response["block"] != nil && response["block"] != "null" {
			blockData, err := json.Marshal(response["block"]) // Convert the block back to JSON
			if err != nil {
				fmt.Printf("Error marshalling block data from %s: %v\n", IP, err)
				continue
			}

			var block Block
			if err := json.Unmarshal(blockData, &block); err != nil {
				fmt.Printf("Error unmarshalling block to struct from %s: %v\n", IP, err)
				continue
			}

			fmt.Println("Block received from miner")
			latestMinersBlock = append(latestMinersBlock, block)

		} else {
			fmt.Printf("No valid block received from miner %s.\n", IP)
		}
	}
	if len(latestMinersBlock) == 0 {
		fmt.Println("No valid blocks received from miners.")
		return
	}

	verfiyMinersLatestBlock(latestMinersBlock)
}

/*
find the hash of each block and then pick one block which has the highest same hash
*/
func verfiyMinersLatestBlock(latestMinerBlock []Block) {
	// Get the hash of each block
	hashes := make(map[string]int)
	for _, block := range latestMinerBlock {
		hash, err := hashStruct(block)
		if err != nil {
			fmt.Printf("Error hashing block: %v\n", err)
			continue
		}
		hashes[hash]++
	}

	// Find the block with the highest same hash
	maxHash := 0
	var maxHashBlock Block
	for _, block := range latestMinerBlock {
		hash, err := hashStruct(block)
		if err != nil {
			fmt.Printf("Error hashing block: %v\n", err)
			continue
		}
		if hashes[hash] > maxHash {
			maxHash = hashes[hash]
			maxHashBlock = block
		}
	}

	if ProofAI.ledger.blocks == nil {
		ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, maxHashBlock)
		InsertBlockInLedgerFile(&maxHashBlock)
		fmt.Println("Ledger Updated Successfully")
		return
	}

	latestBlockIndex := len(ProofAI.ledger.blocks) - 1

	// Compare the block with the ledger
	ledgerBlock := ProofAI.ledger.blocks[latestBlockIndex]
	ledgerHash, err := hashStruct(ledgerBlock)
	if err != nil {
		fmt.Printf("Error hashing ledger block: %v\n", err)
		return
	}
	maxHashBlockHash, err := hashStruct(maxHashBlock)
	if err != nil {
		fmt.Printf("Error hashing max hash block: %v\n", err)
		return
	}

	if ProofAI.ledger.blocks[latestBlockIndex].BlockNum == maxHashBlock.BlockNum {
		// now compare the both hashes
		if ledgerHash != maxHashBlockHash {
			fmt.Println("Ledger block hash does not match with the max hash block")
			ProofAI.ledger.blocks[latestBlockIndex] = maxHashBlock
			fmt.Println("Ledger updated successfully")
			//	InsertBlockInLedgerFile(&maxHashBlock)
		}
	} else {

		if ProofAI.ledger.blocks[latestBlockIndex].BlockNum > maxHashBlock.BlockNum {
			fmt.Println("Ledger is already updated")
			return
		}
		ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, maxHashBlock)
		InsertBlockInLedgerFile(&maxHashBlock)
		fmt.Println("Ledger updated successfully")
		return
	}

	fmt.Println("Ledger is already correct  ")

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
 1. block: block object
    each miner execute each transaction in the block and find own trained model and then verify the block by hash
    If the block is valid, add it to the ledger
    If the block is invalid, mine the block again where it paused
*/
func IncomingBlockVerfication(block *Block) {

	fmt.Println("Incoming Block Verification started")

	if ProofAI.selfMiningDetail.role == "Miner" {
		for _, transaction := range block.Transactions {

			fmt.Printf("Block Transaction nonce: %d, From: %s\n", transaction.Nonce, transaction.From)
			transactionExist := findBlockBy_Nonce_From(transaction.Nonce, transaction.From)
			if !transactionExist {
				fmt.Println("Above Transaction is need to be mined.")
				MineTransaction(&transaction, &ProofAI.selfMiningDetail.CurrentlyMineBlock)
			}
		}
		fmt.Println("Transaction verification completed")
		isValidBlock, err := IsIncomingBlockValid(block, ProofAI.selfMiningDetail.CurrentlyMineBlock.Transactions) // error in this
		if err != nil {
			fmt.Printf("Error during block verification: %v\n", err)
			return
		}

		if isValidBlock {
			return
		}
		fmt.Println("Incoming block is invalid. Mining the block again start .")
		err = PoW(ProofAI.CurrentlyMineBlock, nil)
		if err != nil {
			fmt.Printf("Error during Proof of Work for block: %v\n", err)
			BlockMiningEnd()
			return
		}

		// Finalize block details

		// Log block mining success
		fmt.Println("Block mined successfully. Broadcasting to all miners...")

		// Add to receivedBlock and broadcast
		ProofAI.receivedBlock[ProofAI.CurrentlyMineBlock.TransactionsHash] = true

		broadcastTransaction(&ProofAI.Miners, ProofAI.CurrentlyMineBlock)

		// Log broadcast completion
		fmt.Printf("Block broadcasted successfully at %s.\n", time.Now().Format(time.RFC3339))
	} else {
		fmt.Println("Only verified by POW.")
		ProofAI.CurrentlyMineBlock = block
		ProofAI.selfMiningDetail.CurrentlyMineBlock = *block
	}

	ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, *ProofAI.CurrentlyMineBlock)
	InsertBlockInLedgerFile(ProofAI.CurrentlyMineBlock)
	BlockMiningEnd()
}

/*
IsIncomingBlockValid is function used to verify the incoming block after each transaction in model and attach with this block and then find hash of incoming block and compare with current block hash
 1. block: block object
 2. transactions: list of transactions
    Compute the hash of the incoming block
    Compute the hash of the current block
    If the hashes match, add the block to the ledger
*/
func IsIncomingBlockValid(block *Block, transactions []Transaction) (bool, error) {

	tempSelfMiningBlock := block
	for i, transaction := range transactions {
		tempSelfMiningBlock.Transactions[i] = transaction
	}

	currentBlockHash, err := hashStruct(tempSelfMiningBlock)
	if err != nil {
		fmt.Printf("Error during hash generation: %v\n", err)
		return false, err
	}

	incomingBlockHash, err := hashStruct(block)
	if err != nil {
		fmt.Printf("Error during hash generation for incoming block: %v\n", err)
		return false, err
	}

	if currentBlockHash != incomingBlockHash {
		fmt.Println("Error: Incoming block hash does not match. Incoming block is invalid.")
		return false, nil
	}

	fmt.Println("Incoming block is valid and will now be added to the ledger.")
	ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, *block)
	InsertBlockInLedgerFile(block)
	BlockMiningEnd()

	return true, nil
}

/*
generateBlock is a function to generate a block and mine it
1- transactions: list of transactions
2- wg: wait group object
3- ctx: context object

	Get the previous block hash
	If the previous block hash does not have the required prefix, set it to the genesis block hash
	Set the block number
	Set the proposer ID
	Set the difficulty level
	Set the transactions
	Process the transactions
	Compute the hash of the transactions
	Set the block type
	Set the timestamp
	Perform Proof of Work
*/
func generateBlock(trans_list []Transaction, wg *sync.WaitGroup, ctx context.Context) {

	defer wg.Done()
	ProofAI.selfMiningDetail.interuptStatus = false
	ProofAI.difficultyLevel = ProofAI.selfMiningDetail.powLenght
	updateLedger()

	ProofAI.selfMiningDetail.CurrentlyMineBlock = Block{}
	ProofAI.CurrentlyMineBlock = &ProofAI.selfMiningDetail.CurrentlyMineBlock
	ProofAI.CurrentlyMineBlock.Transactions = trans_list

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
		prev_blockHash = GenesisBlockHash()
		ProofAI.CurrentlyMineBlock.BlockNum = 1
	}

	if !strings.HasPrefix(prev_blockHash, strings.Repeat("0", ProofAI.difficultyLevel)) {
		prev_blockHash = GenesisBlockHash()
	}

	ProofAI.CurrentlyMineBlock.Prev_Hash = prev_blockHash
	ProofAI.CurrentlyMineBlock.ProposerId = ProofAI.selfMiningDetail.pubKeyStr
	ProofAI.CurrentlyMineBlock.Difficulty = ProofAI.difficultyLevel
	ProofAI.currentlyMiningBlockForUser = *ProofAI.CurrentlyMineBlock
	ProofAI.CurrentlyMineBlock.Transactions = nil

	for _, transaction := range trans_list { // Process transactions
		transaction.BlockNum = ProofAI.CurrentlyMineBlock.BlockNum
		MineTransaction(&transaction, ProofAI.CurrentlyMineBlock)
	}

	trans_hash, err := hashStruct(ProofAI.CurrentlyMineBlock.Transactions)
	if err != nil {
		fmt.Printf("Error constructing block hash of transactions: %v\n", err)
		BlockMiningEnd()
		return
	}
	ProofAI.CurrentlyMineBlock.TransactionsHash = trans_hash
	ProofAI.CurrentlyMineBlock.Type = "block"
	ProofAI.CurrentlyMineBlock.TimeStamp = time.Now().Format(time.RFC3339)
	ProofAI.CurrentlyMineBlock.Difficulty = ProofAI.difficultyLevel

	if ctx.Err() != nil {
		ProofAI.selfMiningDetail.interuptStatus = true
		return
	}

	err = PoW(ProofAI.CurrentlyMineBlock, ctx)
	if err != nil {
		fmt.Printf("Error during Proof of Work for block: %v\n", err)
		BlockMiningEnd()
		return
	}

	ProofAI.receivedBlock[ProofAI.CurrentlyMineBlock.TransactionsHash] = true
	broadcastTransaction(&ProofAI.Miners, ProofAI.CurrentlyMineBlock)
	InsertBlockInLedgerFile(ProofAI.CurrentlyMineBlock)
	ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, *ProofAI.CurrentlyMineBlock)
	BlockMiningEnd()
}

/*
BlockMiningEnd is a function to end block mining

	Set the currently mining block to nil
	Set the currently mining block for the user to an empty block
*/
func BlockMiningEnd() {
	ProofAI.CurrentlyMineBlock = nil
	ProofAI.currentlyMiningBlockForUser = Block{}
}

/*
GenesisBlockHash is a function to generate the hash of the genesis block
 1. return: hash of the genesis block ( 0's of length equal to the block length)
*/
func GenesisBlockHash() string {
	requiredPrefix := strings.Repeat("0", ProofAI.selfMiningDetail.blockLength)
	return requiredPrefix
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
    Generate the required prefix
    Generate the salt
    Compute the hash of the block
    Check if the hash has the required prefix
    If the hash has the required prefix, set the salt and return nil
    If the context is canceled, set the interupt status and return an error
*/
func PoW(block *Block, ctx context.Context) error {

	blockSize := uint(102400)
	saltSize := blockSize - uint(unsafe.Sizeof(block))

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

		if strings.HasPrefix(blockHash, requiredPrefix) {
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
			continue
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

	for ProofAI.selfMiningDetail.role != "Miner" {
		time.Sleep(1 * time.Second)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:

			if len(ProofAI.ledger.blocks) > 0 {

				lastBlock := ProofAI.ledger.blocks[len(ProofAI.ledger.blocks)-1]
				lastBlockTime, err := time.Parse(time.RFC3339, lastBlock.TimeStamp)
				if err != nil {
					fmt.Printf("Error parsing time: %v\n", err)
					return
				}
				currentTime := time.Now()
				diff := currentTime.Sub(lastBlockTime)
				if diff.Minutes() > 2 {

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
    Execute the model
    Add the transaction to the block
*/
func MineTransaction(transaction *Transaction, block *Block) {

	pubkey, err := hexToPublicKey(transaction.From)
	if err != nil {
		fmt.Printf("Error getting public key from hex string %v", err)
		return
	}

	signatureValidation, err := verifyTransaction(pubkey, transaction)
	if err != nil {
		fmt.Printf("Error verifying transaction %v", err)
		return
	}
	if !signatureValidation {
		fmt.Printf("Error: Transaction signature is invalid")
		return
	} else {
		fmt.Println("Transaction signature is valid")
	}

	dirPath := filepath.Join(os.TempDir(), ProofAI.modelExecutionDir+time.Now().Format("20060102150405"))

	if _, err := os.Stat(dirPath); err == nil {
		if err := os.RemoveAll(dirPath); err != nil {
			fmt.Printf("Error removing existing directory: %v\n", err)
			return
		}
		fmt.Println("Existing directory removed")
	} else if !os.IsNotExist(err) {
		fmt.Printf("Error checking directory existence: %v\n", err)
		return
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		cleanDir(dirPath)
		return
	}

	modelOutput, transactionLog, err := modelExecution(transaction.Input_dataSet, transaction.Input_model, dirPath)

	transaction.Model_output = modelOutput
	transaction.TransactionLog = transactionLog
	block.Transactions = append(block.Transactions, *transaction)
	cleanDir(dirPath)
}

/*
cleanDir is a function to clean up a directory
 1. dirpath: path of the directory
    Remove the directory
*/
func cleanDir(dirpath string) {
	if err := os.RemoveAll(dirpath); err != nil {
		fmt.Printf("Error removing directory at cleanup: %v\n", err)
		return
	}
	fmt.Println("Temporary directory cleaned up")
}

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
	broadcastTransaction(&ProofAI.Miners, &transaction_)

	return transaction_, nil
}
