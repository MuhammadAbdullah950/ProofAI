package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// write function which write Block struct to file
func InsertBlockInledgerFile(filename string, block* Block) error {

	file , err :=os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err !=nil{
		log.Fatalf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	blockdata , err :=json.Marshal(block)
	if err !=nil{
		log.Fatalf("Error marshalling block: %v", err)
		return err
	}

	file.Write(append(blockdata, '\n'))
	return nil
}


func writeIPTable(filename string, data string) error {
	data += "\n"
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return err
	}
	parts := strings.Split(string(fileContent), "\n")

	for _, part := range parts {
		if part == strings.TrimSpace(data) {
			//	fmt.Println("Data already exists in the file.")
			return nil
		}
	}
	_, err = file.Seek(0, os.SEEK_END)
	if err != nil {
		log.Fatalf("Error seeking to the end of the file: %v", err)
		return err
	}

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
		return err
	}
	return nil
}



func readIPTable(filename string) (string, *ecdsa.PublicKey, error) {

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	IPs := string(data)
	IP_Table := strings.Split(IPs, "\n")
	rand.Seed(time.Now().UnixNano())

	if len(IP_Table) == 1 {
		return "", nil, err
	}

	random_num := rand.Intn(len(IP_Table) - 1)
	minerConn := IP_Table[random_num]
	parts := strings.Split(minerConn, " ")
	pubKey, err := hexToPublicKey(parts[1])
	if err != nil {
		return "", nil, err
	}

	return parts[0], pubKey, err
}

func readFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return "", err

	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return "", err
	}
	//fmt.Println("Data read from file: ", string(data))
	return string(data), nil
}


func getPubKeyofIP(serverURL string,  IP string) ( *ecdsa.PublicKey, error) {

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
        return  nil, nil
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

// func getPubKeyofIP1(filename string, IP string) (*ecdsa.PublicKey, error) {

// 	// Open the file
// 	data, err := readFromFile(filename)

// 	MinerTable := strings.Split(data, "\n")

// 	for _, minerConn := range MinerTable {
// 		parts := strings.Split(minerConn, " ")
// 		if parts[0] == IP {
// 			pubKey, err := hexToPublicKey(parts[1])
// 			if err != nil {
// 				return nil, err
// 			}
// 			return pubKey, err
// 		}
// 	}
// 	return nil, err

// }
