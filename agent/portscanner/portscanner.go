package portscanner

import (
	"fmt"
	"net/netip"
	"os/exec"
	"strings"
	"sync"
)

type Host struct {
	IpAddress  string
	MacAddress string
	Dynamic    bool
}

func RunPortscan(cidr string) []Host {
	ips := getIps(cidr)
	execPings(ips)
	fmt.Println("DONE!")
	return parseArpTable()
}

func getIps(cidr string) []netip.Addr {
	prefix, _ := netip.ParsePrefix(cidr)
	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}
	return ips
}

func execPings(ips []netip.Addr) {
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip_address netip.Addr) {
			defer wg.Done()
			cmd := exec.Command("ping", ip_address.String())
			cmd.Stdout = nil
			cmd.Stderr = nil
			cmd.Run()
		}(ip)
	}
	wg.Wait()
}

func parseArpTable() []Host {
	data, _ := exec.Command("arp", "-a").Output()
	var hosts = []Host{}

	skipHeaders := false
	for _, line := range strings.Split(string(data), "\n") {
		// Skip empty lines
		if len(line) <= 0 {
			continue
		}
		// Skip interfaces lines
		if line[0] != ' ' {
			skipHeaders = true
			continue
		}
		// Skip headers
		if skipHeaders {
			skipHeaders = false
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		var host = Host{}
		host.IpAddress = fields[0]
		host.MacAddress = fields[1]
		if fields[2] == "dynamic" {
			host.Dynamic = true
		} else {
			host.Dynamic = false
		}
		hosts = append(hosts, host)
	}
	return hosts
}
