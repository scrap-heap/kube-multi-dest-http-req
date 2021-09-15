package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {
	flag.Parse()

	if n := flag.NArg(); n != 1 {
		log.Fatalf("supported only one narg, given %d", n)
	}

	// lookup dns names
	ips, err := net.LookupIP(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed lookup command: %s\n", err)
	}

	// extract only v4 addresses
	// uncommet if you get ipv6 addresses in lookup
	// and ipv6 not supported on your hosts

	// allIPs := make([]net.IP, len(ips))
	// copy(allIPs, ips)

	// ips = make([]net.IP, 0)

	// for _, ip := range allIPs {
	// 	if ip.To4() == nil {
	// 		continue
	// 	}

	// 	ips = append(ips, ip)
	// }

	// make request
	body, err := json.Marshal(map[string]int64{
		"request_id": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalf("failed to marshal body bytes: %s\n", err)
	}

	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)

		go func(ip net.IP) {
			defer wg.Done()

			u := url.URL{
				Scheme: "http",
				Host:   ip.String(),
			}

			buf := bytes.NewBuffer(body)

			resp, err := http.Post(u.String(), "", buf)
			if err != nil {
				log.Fatalf("failed to make http request %s: %s", u.String(), err)
			}

			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("failed to read response body %s %s: %s", u.String(), resp.Status, err)
			}

			log.Println(string(body))
		}(ip)
	}

	wg.Wait()
}
