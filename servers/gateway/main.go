package main

import (
	"GroupPoop/servers/gateway/sessions"
	"GroupPoop/servers/gateway/handlers"
	"GroupPoop/servers/gateway/proxy"
	"os"
	"fmt"
	"net/http"
	"net/url"
	"net/http/httputil"
	"log"
	"time"
	"strings"
	"github.com/go-redis/redis"
)

func main() {
	// connect to redis cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISADDR"), // use ENV variable REDISADDR
		Password: "",                     // no password set
		DB:       0,                      // use default DB
	})

	pong, err := rdb.Ping().Result()
	log.Println(pong, err)

	redisStore := sessions.NewRedisStore(rdb, 10*time.Minute)

	signingKey := os.Getenv("SESSIONKEY")

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

	// Receive address(es) for messages microservice(s) and insert into CustomDirector
	bathroomsAddr := os.Getenv("BATHROOMADDR")
	bathroomsAddrs := strings.Split(bathroomsAddr, ",")
	bathroomsAddrURLs := []*url.URL{}
	for i, _ := range bathroomsAddrs {
		bathroomsAddrURL, err := url.Parse(bathroomsAddrs[i])
		if err != nil {
			fmt.Printf("error parsing messages URLs: %v\n", err)
			os.Exit(1)
		}
		bathroomsAddrURLs = append(bathroomsAddrURLs, bathroomsAddrURL)
	}
	bathroomsProxy := &httputil.ReverseProxy{Director: proxy.CustomDirector(bathroomsAddrURLs)}

	mux := http.NewServeMux()

	mux.Handle("/")
	mux.HandleFunc("/", HelloServer)

	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}