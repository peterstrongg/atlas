package main

import (
	"atlas/agent/portscanner"
	"encoding/json"
	"fmt"
)

func main() {
	output := portscanner.RunPortscan("192.168.0.0/24")
	sendData(output)
}

// Sends data to reporting server
func sendData(nmapData []portscanner.Host) {
	jsonData, err := json.Marshal(nmapData)
	if err != nil {
		// TODO: Error handling
	}
	fmt.Println(string(jsonData))
	// TODO: Make REST API request to reporting server
}
