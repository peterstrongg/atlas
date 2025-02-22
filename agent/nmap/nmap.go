package nmap

import (
	"encoding/xml"
	"net"
	"net/netip"
	"os/exec"
)

type NmapStruct struct {
	Hosts []NmapHost `xml:"host"`
}

type NmapHost struct {
	Addresses []NmapAddress `xml:"address"`
}

type NmapAddress struct {
	Address     string `xml:"addr,attr"`
	AddressType string `xml:"addrtype,attr"`
}

func RunScan() NmapStruct {
	// TODO: Pass optional IP and CIDR argument to specify ip range without calling getIpSubnet()
	// Automatic CIDR detection should be secondary to hard code
	ipSubnet := getIpSubnet()
	return parseScan(runNmapScan(ipSubnet))
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

func parseScan(scanOutput []byte) NmapStruct {
	nmapStruct := NmapStruct{}
	xml.Unmarshal(scanOutput, &nmapStruct)
	return nmapStruct
}
