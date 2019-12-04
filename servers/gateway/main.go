package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	if tlsKeyPath == "" {
		// TODO print error to standard out
		fmt.Fprintf(os.Stderr, "error: KeyPath not set %v\n", 300)
		os.Exit(1)
	}

	if tlsCertPath == "" {
		// TODO print error to standard out
		fmt.Fprintf(os.Stderr, "error: CertPath not set %v\n", 300)
		os.Exit(1)
	}

	// Receive address(es) for bathroom microservice(s) and insert into CustomDirector
	bathroomAddr := os.Getenv("BATHROOMADDR")
	bathroomAddrs := strings.Split(bathroomAddr, ",")
	bathroomAddrURLs := []*url.URL{}
	for i, _ := range bathroomAddrs {
		bathroomAddrURL, err := url.Parse(bathroomAddrs[i])
		if err != nil {
			fmt.Printf("error parsing bathroom URLs: %v\n", err)
			os.Exit(1)
		}
		bathroomAddrURLs = append(bathroomAddrURLs, bathroomAddrURL)
	}
	BathroomProx := &httputil.ReverseProxy{Director: CustomDirector(bathroomAddrURLs)}

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloServer)
	mux.Handle("/bathroom", bathroomProxy)
	mux.Handle("/bathroom/", bathroomProxy)
	mux.Handle("user/:userID/review/", bathroomProxy)

	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
