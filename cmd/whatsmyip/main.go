package main

import (
	"fmt"
	"sort"

	whatsmyip "github.com/TechMDW/whatsmyip/pkg"
)

func main() {
	ips := whatsmyip.GetIp()

	if len(ips) == 0 {
		fmt.Println("No IPs found")
		return
	} else if len(ips) == 1 {
		fmt.Println(ips[0].Ipv4)
		return
	}

	sort.Slice(ips, func(i, j int) bool {
		return ips[i].Certainty > ips[j].Certainty
	})

	for _, ip := range ips {
		fmt.Printf("IP: %s | Certainty: %.1f%%\n", ip.Ipv4, ip.Certainty)
	}
}
