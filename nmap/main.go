package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"net/netip"
	"os/exec"
)

type NmapStruct struct {
	NmapHosts []NmapHost `xml:"host"`
}

type NmapHost struct {
	NmapAddresses []NmapAddress `xml:"address"`
}

type NmapAddress struct {
	Address     string `xml:"addr,attr"`
	AddressType string `xml:"addrtype,attr"`
}

func main() {
	ipSubnet := getIpSubnet()
	out := runNmapScan(ipSubnet)
	parseScan(out)
}

func getIpSubnet() netip.Prefix {
	addrs, _ := net.InterfaceAddrs()

	var ipSubnet netip.Prefix
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			// ip, _ := netip.ParseAddr(ipNet.IP.String())
			prefix, _ := netip.ParsePrefix(ipNet.String())
			ipSubnet = prefix.Masked()
		}
	}

	return ipSubnet
}

func runNmapScan(ipSubnet netip.Prefix) []byte {
	// nmap -sP -n -oX - 192.168.0.0/24
	cmd := exec.Command("nmap", "-sP", "-n", "-oX", "-", ipSubnet.String())
	output, _ := cmd.CombinedOutput()
	return output
}

func parseScan(scanOutput []byte) {
	nmapStruct := NmapStruct{}
	xml.Unmarshal(scanOutput, &nmapStruct)
	fmt.Println(nmapStruct)
}
