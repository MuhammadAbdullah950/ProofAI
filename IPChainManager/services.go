package main

/*
  First need to run the IPFS daemon
  And in this file, we are creating a server that will listen on port 8081
  And we have 3 endpoints:
  1. /fetch - POST request to fetch files from IPFS
  2. /machines - GET request to get all machines
  3. /machine - POST request to add a new machine
*/
import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

/*
FileInfo struct is used to store information about a file
*/
type FileInfo struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

/*
Response struct is used to send response to the client
*/
type Response struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Files   []FileInfo `json:"files"`
}

/*
MachineDetail struct is used to store information about a machine
*/
type MachineDetail struct {
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	PubKey    string    `json:"pubKey"`
	Timestamp time.Time `json:"timestamp"`
}

/*
Server struct is used to store information about all machines
*/
type Server struct {
	machines map[string]MachineDetail
	mutex    sync.RWMutex
}

/*
NewServer function is used to create a new Server object
*/
func NewServer() *Server {
	return &Server{
		machines: make(map[string]MachineDetail),
	}
}

/*
handleGetMachines function is used to handle GET request to get all machines
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
handleAddMachine function is used to handle POST request to add a new machine
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(machine)
}

/*
handleRequest function is used to handle POST request to fetch files from IPFS
*/
func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

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
main function is the entry point of the program
*/
func main() {

	server := NewServer()

	http.HandleFunc("/fetch", handleRequest)
	http.HandleFunc("/machines", server.handleGetMachines)
	http.HandleFunc("/machine", server.handleAddMachine)

	fmt.Println("IPChainManager start listening on 0.0.0.0:8081...")
	http.ListenAndServe("0.0.0.0:8081", nil)
}
