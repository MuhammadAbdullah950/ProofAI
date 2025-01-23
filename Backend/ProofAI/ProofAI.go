package main

//			Starting point of the application
//			1. Create a new ProofAI object
//			2. Start the server and listen for incoming requests

// ProofAI is a global variable for session management
var ProofAI *ProofAIFactory

// serviceMachineAdd is the address of the service machine , set before session creating
var serviceMachineAdd string

// createServerAndListen creates a new ProofAI object and starts the server
func main() {

	// create server and listen for incoming requests ( REST API )
	createServerAndListen()
}
