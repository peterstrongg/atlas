package main

import (
	"fmt"
	"net/netip"
	"os/exec"
	"sync"
)

func main() {
	cidr := "192.168.0.0/24"
	ips := getAllIPs(cidr)

	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip_address netip.Addr) {
			defer wg.Done()
			execPing(ip_address)
		}(ip)
	}
	wg.Wait()
	fmt.Println("Done!")
}

func getAllIPs(cidr string) []netip.Addr {
	prefix, _ := netip.ParsePrefix(cidr)
	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}
	return ips
}

func execPing(ip netip.Addr) {
	cmd := exec.Command("ping", ip.String())
	output, _ := cmd.Output()
	fmt.Println(output)
}
