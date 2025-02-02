package main

/*
	In this we create a new session of ProofAI when user login and close the current session when user logs out.
	1-		StartNewSession is a function to start a new session of ProofAI.
	2-		CloseSession is a function to close the current session of ProofAI.
	3-		ProofAIFactory is a struct to create a new ProofAI object.
	4-		NewProofAIFactory is a function to create a new ProofAI object.
	5-		Reset is a function to reset the ProofAI object.
*/

import (
	"context"
	"time"
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
	ProofAI.selfMiningDetail.connListen.Close()
	time.Sleep(2 * time.Second)
	ProofAI.Reset()
	sendServiceLogout()
}

/*
ProofAIFactory is a struct to create a new ProofAI object
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
Purpose of the function is to create a new ProofAIFactory object and initialize it with default values
*/
func NewProofAIFactory() *ProofAIFactory {
	return &ProofAIFactory{
		startPort:           8280,
		connectionPort:      "8090",
		modelExecutionDir:   "TransactonExecution",
		selfMiningDetail:    selfMiner{nonce: 0, role: "Miner", connectionAlive: true, serviceMachineAddr: serviceMachineAdd, readLedger: false},
		memPool:             MemPool{},
		Miners:              []Miner{},
		ledger:              Ledger{},
		CurrentlyMineBlock:  nil,
		receivedTransaction: make(map[string]bool),
		receivedBlock:       make(map[string]bool),
	}
}

/*
Reset resets the ProofAI object, is used to reset the ProofAI object
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
