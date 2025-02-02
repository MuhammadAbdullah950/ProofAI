package main

/*
This file contains the code for the service machine. The service machine is used to provide the service to the miner machines.
The service machine is used to provide the following services:
1. Upload and pin data on IPFS
2. Fetch the files from IPFS CID
3. Add the miner machine
4. Get the list of miner machines
5. Remove the miner machine
6. Check if the miner machine is live or not
7. Get the IPFS CID from the miner machine

*/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

/*
FileInfo struct is used to store the information of the files
*/
type FileInfo struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

/*
Response struct is used to store the response of the request for the files
*/
type Response struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Files   []FileInfo `json:"files"`
}

/*
MachineDetail struct is used to store the information of the miner machine
*/
type MachineDetail struct {
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	PubKey    string    `json:"pubKey"`
	Timestamp time.Time `json:"timestamp"`
}

/*
Server struct is used to store the information of the server
*/
type Server struct {
	machines map[string]MachineDetail
	mutex    sync.RWMutex
}

/*
ChainInfo struct is used to store the information of the chain . And it is used to store the PoW length and proof
*/
type ChainInfo struct {
	PowLen int `json:"powLen"`
	Proof  int `json:"proof"`
}

/*
printTitle function is used to print the title of the service machine
*/
func printTitle(text string) {
	terminalWidth := 120
	padding := (terminalWidth - len(text)) / 2

	boldBlue := "\033[1;34m"
	reset := "\033[0m"

	fmt.Println(strings.Repeat(" ", padding) + boldBlue + text + reset + "\n\n")
}

/*
getRadminIPv4 function is used to get the IPv4 address of the RadminVPN. Reason to use RadminVPN is to connect in same network
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
NewServer function is used to create a new server instance
*/
func NewServer() *Server {
	return &Server{
		machines: make(map[string]MachineDetail),
	}
}

/*
handleGetMachines function is used to handle the get request for the machines
*/
func (s *Server) handleGetMachines(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mutex.RLock()
	machines := make([]MachineDetail, 0, len(s.machines))
	for _, machine := range s.machines {
		machines = append(machines, machine)
	}
	s.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(machines)
}

/*
handleAddMachine function is used to handle the post request for the machines to add the machine
*/
func (s *Server) handleAddMachine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var machine MachineDetail
	if err := json.NewDecoder(r.Body).Decode(&machine); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	machine.Timestamp = time.Now()

	s.mutex.Lock()
	s.machines[machine.IP] = machine
	s.mutex.Unlock()

	fmt.Println("Machine is added")

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chainInfo)
}

/*
handleRequest function is used to handle the request for the files from IPFS CID
*/
func handleRequest(w http.ResponseWriter, r *http.Request) {

	// Get CID from request
	cid := r.FormValue("cid")
	if cid == "" {
		http.Error(w, "CID is required", http.StatusBadRequest)
		return
	}

	// Connect to IPFS
	sh := shell.NewShell("localhost:5001")

	// List all files in directory
	files, err := sh.List(cid)
	if err != nil {
		http.Error(w, "Failed to list IPFS directory", http.StatusInternalServerError)
		return
	}

	// Get content for each file
	var fileInfos []FileInfo
	for _, file := range files {
		reader, err := sh.Cat(file.Hash)
		if err != nil {
			continue
		}

		// Read content
		content := make([]byte, file.Size)
		reader.Read(content)
		reader.Close()

		fileInfos = append(fileInfos, FileInfo{
			Name:    file.Name,
			Hash:    file.Hash,
			Content: string(content),
		})
	}

	// Send response
	response := Response{
		Success: true,
		Message: "Files retrieved successfully",
		Files:   fileInfos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*
IsMinerLive function is used to check if the miner is live or not by pinging the miner machine and remove the offline machines
*/
func (s *Server) IsMinerLive() {
	for {
		s.mutex.RLock()
		// Collect IPs of offline miners to remove
		var offlineMachines []string

		for ip, machine := range s.machines {
			resp, err := http.Get("http://" + machine.IP + ":8079/ping")

			if err != nil {
				// Check if the error indicates a connection refusal
				if isConnectionRefused(err) {
					//	fmt.Printf("Error pinging miner %s: %v\n", machine.IP, err)
					offlineMachines = append(offlineMachines, ip) // Mark as offline
				} else {
					//fmt.Printf("Unexpected error for miner %s: %v\n", machine.IP, err)
				}
				continue
			}

			// Check the response status code
			resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound { // 404 indicates the miner is live
				//fmt.Printf("Miner %s returned status: 404 (Miner is live)\n", machine.IP)
			} else {
				//fmt.Printf("Miner %s returned status: %d (Unexpected)\n", machine.IP, resp.StatusCode)
			}
		}
		s.mutex.RUnlock()

		// Remove offline machines outside the read lock
		if len(offlineMachines) > 0 {
			s.mutex.Lock()
			for _, ip := range offlineMachines {
				delete(s.machines, ip)
				fmt.Printf("Removed offline miner: %s\n", ip)
			}
			s.mutex.Unlock()
		}

		time.Sleep(1 * time.Second) // Wait 1 second before the next check
	}
}

/*
isConnectionRefused function is used to check if the connection is refused
*/
func isConnectionRefused(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*net.OpError); ok {
			return strings.Contains(syscallErr.Error(), "connection refused")
		}
	}
	return strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "actively refused")
}

/*
handleuploadAndPinData function is used to handle the upload and pin data on IPFS request
*/
func handleuploadAndPinData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(5 * 1024 << 20) // Max upload size:
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files provided in the request", http.StatusBadRequest)
		return
	}

	tempDir, err := os.MkdirTemp("", "ipfs-upload-*")
	if err != nil {
		http.Error(w, "Failed to create temporary directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up after use

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open file %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		destPath := filepath.Join(tempDir, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create file %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
			return
		}
	}

	sh := shell.NewShell("localhost:5001")
	if !sh.IsUp() {
		http.Error(w, "IPFS node is not running on localhost:5001", http.StatusInternalServerError)
		return
	}

	cid, err := sh.AddDir(tempDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload directory to IPFS: %v", err), http.StatusInternalServerError)
		return
	}
	response := Response{
		Success: true,
		Message: cid,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*
isIPFSRuning function is used to check if the IPFS is running on localhost:5001
*/
func isIPFSRuning() bool {
	sh := shell.NewShell("localhost:5001")
	return sh.IsUp()
}

/*
waitToCloseWindow function is used to wait for the user to close the window by pressing any key
*/
func waitToCloseWindow() {
	fmt.Println("")
	fmt.Println("Press any key to close the window")
	var input string
	fmt.Scanln(&input)
}

/*
handleLogout function is used to handle the logout request from the machine
*/
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr

	// Loop through machines to find the matching IP
	for _, machine := range s.machines {
		if machine.IP == strings.Split(addr, ":")[0] {
			s.mutex.Lock()
			delete(s.machines, machine.IP) // Remove the machine entry
			s.mutex.Unlock()

			// Log the removal and send a success response
			fmt.Printf("Machine with IP %s is removed successfully.\n", machine.IP)
			w.WriteHeader(http.StatusOK) // Send HTTP 200 OK response
			w.Write([]byte("Logout successful, machine removed."))
			return
		}
	}

	// If no machine is found, log and send an error response
	fmt.Printf("Logout request from %s, but machine not found.\n", addr)
	w.WriteHeader(http.StatusNotFound) // Send HTTP 404 Not Found response
	w.Write([]byte("Logout failed, machine not found."))
}

// Global variable to store the chain information (PoW length and proof)
var chainInfo ChainInfo

/*
main function is the entry point of the program
Precondition: The IPFS node should be running on localhost:5001
Postcondition: The service machine is started and listening on the RadminVPN IP address
Creates a new server instance and starts the service machine
*/

func main() {

	text := "ProofAI Service Machine is starting..."
	printTitle(text)

	IP, err := getRadminIPv4()
	if err != nil {
		log.Printf("Failed to get RadminVPN IPv4 address:  %v", err)
		waitToCloseWindow()
		return
	}

	if !isIPFSRuning() {
		log.Printf("IPFS node is not running on localhost:5001")
		waitToCloseWindow()
		return
	}

	server := NewServer()
	go server.IsMinerLive()

	http.HandleFunc("/fetch", handleRequest)
	http.HandleFunc("/upload", handleuploadAndPinData)
	http.HandleFunc("/logout", server.handleLogout)
	http.HandleFunc("/machines", server.handleGetMachines)
	http.HandleFunc("/machine", server.handleAddMachine)

	fmt.Printf("Enter the blockHash Size       : ")
	fmt.Scanln(&chainInfo.PowLen)

	fmt.Printf("Enter the Proof of Work length : ")
	fmt.Scanln(&chainInfo.Proof)

	fmt.Println("\n\nService Machine Address  =  ", IP+":8050 \n\n")
	if err := http.ListenAndServe(IP+":8050", nil); err != nil {
		log.Printf("Failed to start Service Machine : %v", err)
		waitToCloseWindow()
		return
	}

}
