package main

import (
	"context"
)

/*
Starting point of the application
 1. Create a new ProofAI object
 2. Start the server and listen for incoming requests
*/
func StartNewSession() {
	ProofAI = NewProofAIFactory()
	ctx, _ := context.WithCancel(context.Background())
	go BlockMining(ctx)
}

/*
CloseSession closes the current session
logic to close the session
 1. Set connectionAlive to false
 2. Close all the connections
 3. Reset the ProofAI object
*/
func CloseSession() {
	ProofAI.selfMiningDetail.connectionAlive = false
	for _, miner := range ProofAI.Miners {
		miner.conn.Close()
	}
	ProofAI.Reset()
}

/*
ProofAIFactory is a factory class to create a new ProofAI object
 1. startPort: starting port for the server
 2. difficultyLevel: difficulty level for mining
 3. connectionPort: port for connection
 4. modelExecutionDir: directory for model execution
 5. IPTable: file to store IP addresses
 6. selfMiningDetail: details of the miner
 7. memPool: memory pool to store transactions
 8. Miners: list of miners
 9. ledger: ledger to store transactions
 10. CurrentlyMineBlock: block currently being mined
 11. receivedTransaction: map to store received transactions
 12. receivedBlock: map to store received blocks
 13. currentlyMiningBlockForUser: block currently being mined for user
*/
type ProofAIFactory struct {
	startPort                   int
	difficultyLevel             int
	connectionPort              string
	modelExecutionDir           string
	IPTable                     string
	selfMiningDetail            selfMiner
	memPool                     MemPool
	Miners                      []Miner
	ledger                      Ledger
	CurrentlyMineBlock          *Block
	receivedTransaction         map[string]bool
	receivedBlock               map[string]bool
	currentlyMiningBlockForUser Block
}

/*
NewProofAIFactory creates a new ProofAIFactory object
*/

func NewProofAIFactory() *ProofAIFactory {
	return &ProofAIFactory{
		startPort:           8081,
		difficultyLevel:     5,
		connectionPort:      "8090",
		modelExecutionDir:   "TransactonExecution",
		IPTable:             "MachineIPTable.txt",
		selfMiningDetail:    selfMiner{nonce: 0, role: "Miner", connectionAlive: true, serviceMachineAddr: serviceMachineAdd},
		memPool:             MemPool{},
		Miners:              []Miner{},
		ledger:              Ledger{},
		CurrentlyMineBlock:  nil,
		receivedTransaction: make(map[string]bool),
		receivedBlock:       make(map[string]bool),
	}
}

/*
Reset resets the ProofAI object
*/
func (bf *ProofAIFactory) Reset() {
	bf.selfMiningDetail = selfMiner{nonce: 0, role: "Miner"}
	bf.memPool = MemPool{}
	bf.Miners = []Miner{}
	bf.ledger = Ledger{}
	bf.CurrentlyMineBlock = nil
	bf.receivedTransaction = make(map[string]bool)
	bf.receivedBlock = make(map[string]bool)
}
