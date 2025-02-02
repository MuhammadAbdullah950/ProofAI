package main

/*
	In this file we manage the connection with the service machine and the miners.
	New miners are registered with the service machine. And then miners are connected to each other.
	Every miner listens for incoming connections and establishes a communication connection with other miners.
	And after connection is established, miners create a new communication port for communication with each other.
	1. MachineDetail: struct to store the details of a machine
	2. ChainInfo: struct to store the chain information
	3. getRadminIPv4: function to get the IPv4 address of the Radmin VPN
	4. getRandomMiner: function to get a random miner from the service machine
	5. registerMiner: function to register a miner with the service machine
	6. establishConnection: function to establish a connection with the service machine
	7. connectToMiner: function to connect to a miner
	8. writeToConnection: function to write to a connection
	9. handleConnection: function to handle a connection

*/

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
MachineDetail is a struct to store the details of a machine
*/
type MachineDetail struct {
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	PubKey    string    `json:"pubKey"`
	Timestamp time.Time `json:"timestamp"`
}

/*
ChainInfo is a struct to store the chain information
*/
type ChainInfo struct {
	PowLen int `json:"powLen"`
	Proof  int `json:"proof"`
}

/*
getRadminIPv4 is a function to get the IPv4 address of the Radmin VPN
 1. Get the network interfaces
 2. Check if the interface name contains "Radmin"
 3. Get the IP address
 4. Return the IPv4 address
    And it use to make machine publicaly available
*/
func getRadminIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %v", err)
	}

	for _, iface := range interfaces {

		if strings.Contains(iface.Name, "Radmin") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}

				ipv4 := ipNet.IP.To4()
				if ipv4 != nil {
					return ipv4.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("RadminVPN IPv4 address not found")
}

/*
getRandomMiner is a function to get a random miner from the service machine
 1. Send a GET request to the service machine to get the list of miners
 2. Decode the response
 3. Select a random miner , initialy it will be one of the miner from the list
 4. Convert the public key string to *ecdsa.PublicKey
 5. Return the IP address, public key, and error
*/
func getRandomMiner(serverURL string) (string, *ecdsa.PublicKey, error) {
	resp, err := http.Get(serverURL + "/machines")
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch miners: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var machines []MachineDetail
	if err := json.NewDecoder(resp.Body).Decode(&machines); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(machines) == 0 {
		return "", nil, nil
	}

	machine := machines[rand.Intn(len(machines))]
	pubKey, err := hexToPublicKey(machine.PubKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert public key: %v", err)
	}

	return fmt.Sprintf("%s:%s", machine.IP, machine.Port), pubKey, nil
}

/*
registerMiner is a function to register a miner with the service machine
1. Create a MachineDetail object
2. Marshal the object to JSON
3. Send a POST request to the service machine to register the miner
4. Check the response status code
5. Return an error if any
*/
func registerMiner(serverURL, ip, port, pubKeyStr string) error {
	machine := MachineDetail{
		IP:     ip,
		Port:   port,
		PubKey: pubKeyStr,
	}

	jsonData, err := json.Marshal(machine)
	if err != nil {
		return fmt.Errorf("failed to marshal machine data: %v", err)
	}

	resp, err := http.Post(serverURL+"/machine", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to register miner: %v", err)
	}

	var chainInfo ChainInfo
	err = json.NewDecoder(resp.Body).Decode(&chainInfo)
	if err != nil {
		return fmt.Errorf("failed to decode response of Service Machine to Set ChainInfo : %v", err)
	}

	ProofAI.selfMiningDetail.blockLength = chainInfo.PowLen
	ProofAI.selfMiningDetail.powLenght = chainInfo.Proof

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	if !ProofAI.selfMiningDetail.readLedger {
		ReadAndWriteMemoryTransaction()
		ProofAI.selfMiningDetail.readLedger = true
	}
	return nil
}

/*
establishConnection is a function to establish a connection with the service machine
1. Get the machine IP
2. Get a random miner
3. Connect to the miner
4. Listen for incoming connections
5. Register the miner with the service machine
6. Accept incoming connections
7. Establish communication connection
8. Read transactions
9. Write to connection
*/
func establishConnection(port string) {

	machineIP, err := getRadminIPv4()
	if err != nil {
		fmt.Printf("Error getting machine IP: %v\n", err)
		return
	}
	machineIP = strings.ReplaceAll(machineIP, "\n", "")

	serviceMachineURl := "http://" + ProofAI.selfMiningDetail.serviceMachineAddr

	baseMiner, minerPubkey, err := getRandomMiner(serviceMachineURl)
	if err != nil {
		log.Printf("Error reading IPTable: %v\n", err)
		return
	}

	if len(baseMiner) != 0 {
		connectToMiner(baseMiner, minerPubkey)
	}

	IP := "0.0.0.0"
	ln, err := net.Listen("tcp", IP+":"+port)
	ProofAI.selfMiningDetail.connListen = ln

	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer ln.Close()

	for {

		publicIP := machineIP
		publicIP = strings.ReplaceAll(publicIP, "\n", "")
		if err != nil {
			fmt.Printf("Error getting public IP %v\n", err)
			return
		}

		err = registerMiner(serviceMachineURl, publicIP, port, ProofAI.selfMiningDetail.pubKeyStr)
		if err != nil {
			return
		}

		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			return
		} else {

			fmt.Printf("\n\n")
			miner := newMiner(conn)

			addr := miner.conn.RemoteAddr().String()
			parts := strings.Split(addr, ":")

			minerPubkey, err := getPubKeyofIP(serviceMachineURl, parts[0])
			miner.pubKey = minerPubkey

			parts[1] = strconv.Itoa(ProofAI.startPort)
			parts[0] = "0.0.0.0"

			commLn, err := net.Listen("tcp", parts[0]+":"+parts[1])
			if err != nil {
				log.Printf("Error creating communication port for : %s at port: %s \n ", parts[0], parts[1])
			} else {

				go func() {
					commConn, err := commLn.Accept()
					if err != nil {
						log.Printf("Error establishing communication connection with: %s at port: %s ", parts[0], parts[1])
						return
					} else {
						miner.conn = commConn
						ProofAI.Miners = append(ProofAI.Miners, *miner)
						go readTransaction(miner)
						ProofAI.startPort++
					}
				}()
				writeToConnection(miner, parts[1])
			}
		}
	}

}

/*
connectToMiner is a function to connect to a miner
1. Dial the connection
2. Read the communication port
3. Parse the address for communication connection
4. Establish the communication connection
5. Read transactions
*/
func connectToMiner(baseMiner string, minerPubkey *ecdsa.PublicKey) {

	conn, err := net.Dial("tcp", baseMiner)

	if err != nil {
		log.Printf("Error connecting to %s: %v\n", baseMiner, err)
		return
	}

	fmt.Println("Connection established with", baseMiner)
	miner := newMiner(conn)
	miner.pubKey = minerPubkey

	communicationPort, err := miner.read.ReadString('\n')
	if err != nil {
		log.Printf("Error reading from %s: %v", miner.conn.RemoteAddr(), err)
		return
	}
	communicationPort = strings.TrimSpace(communicationPort)

	parts := strings.Split(baseMiner, ":")
	commAddress := fmt.Sprintf("%s:%s", parts[0], communicationPort)
	fmt.Println(commAddress)

	conn, err = net.Dial("tcp", commAddress)
	if err != nil {
		log.Printf("Error establishing communication connection with %s at port %s: %v\n",
			parts[0], communicationPort, err)
		return
	}

	miner.conn = conn
	fmt.Printf("Connection established: %s for communication on port %v\n",
		conn.RemoteAddr(), communicationPort)

	go readTransaction(miner)
	ProofAI.Miners = append(ProofAI.Miners, *miner)
}

/*
writeToConnection is a function to write to a connection
1. Write the message to the connection
2. Log any errors
*/
func writeToConnection(miner *Miner, message string) {
	message = strings.TrimSpace(message) + "\n"
	_, err := miner.conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error writing to client %s: %v \n", miner.conn.RemoteAddr(), err)
		return
	}
}

/*
handleConnection is a function to handle a connection
1. Close the connection when the function returns
2. Create a new miner
3. Get the IP address of the miner
4. Get the public key of the miner
5. Set the public key of the miner
6. Create a communication port
7. Accept incoming connections
8. Establish a communication connection
9. Read transactions
10. Append the miner to the list of miners
11. Increment the starting port
*/
func handleConnection(conn net.Conn) {
	defer conn.Close()

	miner := newMiner(conn)
	addr := miner.conn.RemoteAddr().String()
	parts := strings.Split(addr, ":")

	serviceMachineURl := "http://" + ProofAI.selfMiningDetail.serviceMachineAddr

	minerPubkey, err := getPubKeyofIP(serviceMachineURl, parts[0])
	if err != nil {
		log.Printf("Error getting public key: %v\n", err)
		return
	}
	miner.pubKey = minerPubkey

	parts[1] = strconv.Itoa(ProofAI.startPort)
	parts[0] = "0.0.0.0"

	commLn, err := net.Listen("tcp", parts[0]+":"+parts[1])
	if err != nil {
		log.Printf("Error creating communication port for: %s at port: %s\n", parts[0], parts[1])
		return
	}
	defer commLn.Close()

	commConn, err := commLn.Accept()
	if err != nil {
		log.Printf("Error establishing communication connection with: %s at port: %s\n", parts[0], parts[1])
		return
	}
	defer commConn.Close()

	miner.conn = commConn
	ProofAI.Miners = append(ProofAI.Miners, *miner)

	go readTransaction(miner)
	ProofAI.startPort++
}
