package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

/*
ReadAndWriteMemoryTransaction is a function to read and write the memory transactions
main logic:
 1. Create a file with the name Transaction_<powLenght>_<blockLength>.json
 2. If the file does not exist, create a new file
 3. Read the blocks from the file
 4. Append the blocks to the ledger
*/
func ReadAndWriteMemoryTransaction() {

	blockLengthStr := strconv.Itoa(ProofAI.selfMiningDetail.blockLength)
	powStr := strconv.Itoa(ProofAI.selfMiningDetail.powLenght)
	file := "Transaction_" + powStr + "_" + blockLengthStr + ".json"
	ProofAI.selfMiningDetail.LedgerFile = file

	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("Error creating file: %v\n", err)
		}
		f.Close()
		return
	}

	ProofAI.selfMiningDetail.transactionFile = file

	blocks, err := ReadBlocksFromLedgerFile(file)
	if err != nil {
		log.Fatalf("Error reading transactions from file: %v\n", err)
	}

	for _, block := range blocks {
		ProofAI.ledger.blocks = append(ProofAI.ledger.blocks, *block)
	}
}

/*
InsertBlockInLedgerFile is a function to insert the block in the ledger file
*/
func InsertBlockInLedgerFile(block *Block) error {

	filePath := ProofAI.selfMiningDetail.LedgerFile

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	blockData, err := json.Marshal(block)
	if err != nil {
		log.Fatalf("Error marshalling block: %v", err)
		return err
	}

	if _, err := file.Write(append(blockData, '\n')); err != nil {
		log.Fatalf("Error writing block data to file: %v", err)
		return err
	}
	return nil
}

/*
ReadBlocksFromLedgerFile is a function to read the blocks from the ledger file
*/
func ReadBlocksFromLedgerFile(filename string) ([]*Block, error) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var blocks []*Block

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading file: %v", err)
			return nil, err
		}

		var block Block
		err = json.Unmarshal([]byte(line), &block)
		if err != nil {
			log.Printf("Error unmarshalling block: %v", err)
			return nil, err
		}
		blocks = append(blocks, &block)
	}
	return blocks, nil
}

/*
getPubKeyofIP is a function to get the public key of the IP address from the server URL provided
*/
func getPubKeyofIP(serverURL string, IP string) (*ecdsa.PublicKey, error) {

	resp, err := http.Get(serverURL + "/machines")
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var machines []MachineDetail
	if err := json.NewDecoder(resp.Body).Decode(&machines); err != nil {
		return nil, nil
	}

	if len(machines) == 0 {
		return nil, nil
	}

	for _, machine := range machines {
		if machine.IP == IP {
			fmt.Println(machine.IP)
			pubKey, err := hexToPublicKey(machine.PubKey)
			if err != nil {
				return nil, err
			}
			return pubKey, err
		}
	}
	return nil, nil
}
