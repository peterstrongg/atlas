package main

import (
	"atlas/agent/nmap"
	"fmt"
)

func main() {
	output := nmap.RunScan()
	fmt.Println(output)
}
