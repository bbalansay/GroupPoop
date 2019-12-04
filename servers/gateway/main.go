package main

import (
	"GroupPoop/servers/gateway/sessions"
	"GroupPoop/servers/gateway/middleware"
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

	// Receive address(es) for bathrooms microservice(s) and insert into CustomDirector
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

	// Receive address(es) for auth microservice(s) and insert into CustomDirector
	authAddr := os.Getenv("AUTHADDR")
	authAddrs := strings.Split(authAddr, ",")
	authAddrURLs := []*url.URL{}
	for i, _ := range authAddrs {
		authAddrURL, err := url.Parse(authAddrs[i])
		if err != nil {
			fmt.Printf("error parsing messages URLs: %v\n", err)
			os.Exit(1)
		}
		authAddrURLs = append(authAddrURLs, authAddrURL)
	}
	authProxy := &httputil.ReverseProxy{Director: proxy.CustomDirector(authAddrURLs)}

	// Receive address(es) for useres microservice(s) and insert into CustomDirector
	usersAddr := os.Getenv("USERSADDR")
	usersAddrs := strings.Split(usersAddr, ",")
	usersAddrURLs := []*url.URL{}
	for i, _ := range usersAddrs {
		usersAddrURL, err := url.Parse(usersAddrs[i])
		if err != nil {
			fmt.Printf("error parsing messages URLs: %v\n", err)
			os.Exit(1)
		}
		usersAddrURLs = append(usersAddrURLs, usersAddrURL)
	}
	usersProxy := &httputil.ReverseProxy{Director: proxy.CustomDirector(usersAddrURLs)}

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloServer)
	mux.Handle("/login", authProxy)
	mux.Handle("/login/", authProxy)
	mux.Handle("/user", usersProxy)
	mux.Handle("/user/", usersProxy)
	mux.Handle("/bathroom", bathroomsProxy)
	mux.Handle("/bathroom/", bathroomsProxy)
	mux.Handle("user/:userID/review/", bathroomsProxy)

	wrappedMux := middleware.NewEnsureCORS(middleware.NewEnsureAuth(mux, signingKey, redisStore))
	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
