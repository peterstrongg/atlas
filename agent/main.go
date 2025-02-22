package main

import (
	"atlas/agent/nmap"
	"encoding/json"
	"fmt"
)

func main() {
	output := nmap.RunScan()
	sendData(output)
}

// Sends data to reporting server
func sendData(nmapData nmap.NmapStruct) {
	jsonData, err := json.Marshal(nmapData)
	if err != nil {
		// TODO: Error handling
	}
	fmt.Println(string(jsonData))
}
