package main

import (
	"fmt"
	"net"
	"net/netip"
	"os/exec"
)

func main() {
	ipSubnet := getIpSubnet()
	out := runNmapScan(ipSubnet)
	fmt.Println(out)
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

func runNmapScan(ipSubnet netip.Prefix) string {
	// nmap -sP -n -oX - 192.168.0.0/24
	cmd := exec.Command("nmap", "-sP", "-n", "-oX", "-", ipSubnet.String())
	output, _ := cmd.CombinedOutput()
	return string(output)
}
