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
	var ipsV4 []net.IP
	for _, ip := range ips {
		if ip.To4() == nil {
			continue
		}

		ipsV4 = append(ipsV4, ip)
	}

	// make request
	body, err := json.Marshal(map[string]int64{
		"request_id": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalf("failed to marshal body bytes: %s\n", err)
	}

	var wg sync.WaitGroup
	for _, ip := range ipsV4 {
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
