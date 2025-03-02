package main

import (
	"atlas/agent/pingscan"
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	ReportingUrl string `yaml:"reporting_url"`
}

func main() {
	config := loadConfig("config.yaml")
	output := pingscan.RunPingScan("192.168.0.0/24")
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
func sendData(nmapData []pingscan.Host, reportingUrl string) {
	jsonData, err := json.Marshal(nmapData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))

	// req, err := http.NewRequest("POST", reportingUrl, bytes.NewBuffer(jsonData))
	// req.Header.Set("Content-Type", "application/json")
	// if err != nil {
	// 	panic(err)
	// }

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
}
