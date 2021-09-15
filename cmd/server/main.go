package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	if n := flag.NArg(); n != 1 {
		log.Fatalf("supported only one narg, given %d", n)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read request body: %s", err)
		}

		log.Printf(string(body))
	})

	log.Fatal(http.ListenAndServe(flag.Arg(0), nil))
}
