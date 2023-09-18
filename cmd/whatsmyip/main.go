package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sort"

	whatsmyip "github.com/TechMDW/whatsmyip/pkg"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Output in JSON format")
	jsonRawOutput := flag.Bool("json-raw", false, "Output in JSON format without indentation")
	flag.Parse()

	ips := whatsmyip.GetIp()

	if *jsonRawOutput {
		json, err := json.Marshal(ips)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(json))
		return
	}

	if *jsonOutput {
		json, err := json.MarshalIndent(ips, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(json))
		return
	}

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
