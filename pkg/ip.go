package whatsmyip

import (
	"context"
	"sync"
	"time"
)

type Ip struct {
	// Ip is the IP address.
	Ip string

	// Certainty is the percentage of websites that returned the same IP.
	//
	// This will be 100 if all websites returned the same IP.
	Certainty float64

	// RequestInfo stores the http information for all websites that returned the same IP.
	//
	// This will mainly be used for debugging.
	RequestInfo []IpFetch
}

// Retrieves the IP address by fetching the IP from a list of websites.
func GetIp() (data []Ip) {
	ch := make(chan bool, 6)

	ips := make([]IpFetch, len(Websites))

	wg := sync.WaitGroup{}

	for i, website := range Websites {
		wg.Add(1)
		go func(website string, i int) {
			defer wg.Done()

			ch <- true
			defer func() { <-ch }()

			context, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

			go func() {
				startTime := time.Now()

				// Quick and dirty solution, works for now.
				for {
					select {
					case <-context.Done():
						return
					default:
						found := 0
						for _, ip := range ips {
							if ip.Ip != "" {
								found++
							}
						}

						if time.Since(startTime) > 5*time.Second && found > 4 {
							cancel()
							return
						}
					}
				}
			}()

			ip, err := IpWebScraper(context, website)

			if err != nil {
				// TODO: Figure out how to handle/display errors
				return
			}

			ips[i] = ip
		}(website, i)
	}

	wg.Wait()

	// Count the number of times each IP appears
	ipCount := make(map[string]int)
	ipEmpty := 0
	for _, ip := range ips {
		if ip.Ip == "" {
			ipEmpty++
			continue
		}

		ipCount[ip.Ip]++
	}

	// If all the IPs are the same, return that IP
	if len(ipCount) == 1 {
		for ip := range ipCount {
			if ip != "" {
				httpInfo := getRequestInfoFromIp(ip, &ips)
				data = append(data, Ip{Ip: ip, Certainty: 100, RequestInfo: httpInfo})
				return
			}
		}
	}

	// If we get here, the IPs are not all the same
	for ip, count := range ipCount {
		certainty := float64(count*100) / float64(len(Websites)-ipEmpty)

		httpInfo := getRequestInfoFromIp(ip, &ips)
		data = append(data, Ip{Ip: ip, Certainty: certainty, RequestInfo: httpInfo})
	}

	return
}

// Used to get the http information for a specific IP.
//
// Returns the http information for all websites that returned the passed in IP.
func getRequestInfoFromIp(ip string, ips *[]IpFetch) (data []IpFetch) {
	for _, website := range *ips {
		if website.Ip == ip {
			data = append(data, website)
		}
	}

	return
}
