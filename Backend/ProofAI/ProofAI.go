package main

//			Starting point of the application
//			1. Create a new ProofAIFactory object
//			2. Start the server and listen for incoming requests

// ProofAI is a global variable for session management
var ProofAI *ProofAIFactory

// serviceMachineAdd is the address of the service machine , set before session creating
var serviceMachineAdd string

// createServerAndListen creates a new ProofAIFactory object and starts the server
func main() {
	// Start the external world server
	go createServerAndListenExternelWorld()

	// Start the REST API server
	go createServerAndListen()

	// Keep the main function alive to allow goroutines to run
	select {}
}
