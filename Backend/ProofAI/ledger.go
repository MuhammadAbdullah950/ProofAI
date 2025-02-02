package main

/*
	In this file we store the ledger details.
	1. Ledger: struct to store the ledger details
	2. Block: struct to store the block details
	3. Transaction: struct to store the transaction details
*/

/*
Ledger is a struct to store the ledger details
-Blocks: list of blocks in the ledger
*/
type Ledger struct {
	blocks []Block
}

/*
Block is a struct to store the block details
*/
type Block struct {
	Transactions     []Transaction `json:"transactions"`
	Prev_Hash        string        `json:"prev_Hash"`
	ProposerId       string        `json:"proposerId"`
	BlockNum         int           `json:"blockNum"`
	TimeStamp        string        `json:"timeStamp"`
	TransactionsHash string        `json:"transactionsHash"`
	Salt             string        `json:"salt"`
	Difficulty       int           `json:"difficulty"`
	Type             string        `json:"type"`
}

/*
Transaction is a struct to store the transaction details
*/
type Transaction struct {
	From           string `json:"from"`
	Nonce          int    `json:"nonce"`
	Input_dataSet  string `json:"input_dataSet"`
	Input_model    string `json:"input_model"`
	Model_output   []byte `json:"model_output"`
	TransactionLog []byte `json:"transactionLog"`
	BlockNum       int    `json:"blockNum"`
	Signature      string `json:"signature"`
	Type           string `json:"type"`
}
