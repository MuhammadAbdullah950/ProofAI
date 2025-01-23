package main

/*
	In this file we have the structs used in the application.
	1. Ledger: struct to store the ledger details
	2. Block: struct to store the block details
	3. Transaction: struct to store the transaction details
*/

/*
Ledger is a struct to store the ledger details
 1. Blocks: list of blocks in the ledger
*/
type Ledger struct {
	blocks []Block
}

/*
Block is a struct to store the block details
 1. Transactions: list of transactions in the block
 2. Prev_Hash: hash of the previous block
 3. ProposerId: id of the proposer
 4. BlockNum: block number
 5. EventEmit: event emitted
 6. TimeStamp: timestamp of the block
 7. TransactionsHash: hash of the transactions
 8. Salt: salt of the block
 9. Difficulty: difficulty level of the block
 10. Proof: proof of the block
 11. Type: type of the block
*/
type Block struct {
	Transactions     []Transaction `json:"transactions"`
	Prev_Hash        string        `json:"prev_Hash"`
	ProposerId       string        `json:"proposerId"`
	BlockNum         int           `json:"blockNum"`
	EventEmit        string        `json:"eventEmit"`
	TimeStamp        string        `json:"timeStamp"`
	TransactionsHash string        `json:"transactionsHash"`
	Salt             string        `json:"salt"`
	Difficulty       int           `json:"difficulty"`
	Proof            string        `json:"proof"`
	Type             string        `json:"type"`
}

/*
Transaction is a struct to store the transaction details
 1. From: address of the sender
 2. Nonce: nonce of the transaction
 3. Input_dataSet: input data set for the model
 4. Input_model: input model for the model
 5. Model_output: output of the model
 6. BlockId: block id of the transaction
 7. Signature: signature of the transaction
 8. ModelFile: model file
 9. Type: type of the transaction
*/
type Transaction struct {
	From          string `json:"from"`
	Nonce         int    `json:"nonce"`
	Input_dataSet string `json:"input_dataSet"`
	Input_model   string `json:"input_model"`
	Model_output  []byte `json:"model_output"`
	BlockId       string `json:"blockId"`
	Signature     string `json:"signature"`
	ModelFile     string `json:"modelFile"`
	Type          string `json:"type"`
}
