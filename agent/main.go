package main

import (
	"atlas/agent/portscanner"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	ReportingUrl string `yaml:"reporting_url"`
}

func main() {
	config := loadConfig("config.yaml")
	output := portscanner.RunPortscan("192.168.0.0/24")
	sendData(output, config.ReportingUrl)
}

func loadConfig(path string) yamlConfig {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config yamlConfig
	yaml.Unmarshal(f, &config)
	fmt.Println(config)
	return config
}

// Sends data to reporting server
func sendData(nmapData []portscanner.Host, reportingUrl string) {
	jsonData, err := json.Marshal(nmapData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", reportingUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
