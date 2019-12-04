package main

import (
	"assignments-zhouyifan0904/servers/gateway/handlers"
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type Director func(r *http.Request)

func CustomDirector(targets []*url.URL) Director {
	counter := 0

	return func(r *http.Request) {
		targ := targets[counter%len(targets)]
		counter++
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
	}
}

//main is the main entry point for the server
func main() {
	// connect to redis cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISADDR"), // use ENV variable REDISADDR
		Password: "",                     // no password set
		DB:       0,                      // use default DB
	})

	pong, err := rdb.Ping().Result()
	fmt.Println(pong, err)

	redisStore := sessions.NewRedisStore(rdb, 10*time.Minute)

	// connect to mysql database
	// Bradley: set database password as an environment variable
	// and connect to MySQL database
	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
	} else {
		fmt.Printf("successfully connected!\n")
	}

	// connect to mq server
	conn, err := amqp.Dial(os.Getenv("RABBITADDR"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	for err != nil {
		// keep trying
		fmt.Printf("Error: %v, trying again\n", err)
		time.Sleep(1 * time.Second)
		conn, err = amqp.Dial(os.Getenv("RABBITADDR"))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"NewMessageApis", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		fmt.Printf("Failed to declare a queue %v\n", err)
	}

	userStore := users.NewMySQLStore(db)
	signingKey := os.Getenv("SESSIONKEY")

	ctx := handlers.NewHandlerContext(signingKey, redisStore, userStore, handlers.SocketStore{
		Connections: make(map[int64]*websocket.Conn),
	})

	// start goroutine broadcasting rabbit mq message
	ctx.Broadcast(ch, q)

	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80" */
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production

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

	// Receive address(es) for summary microservice(s) and insert into CustomDirector
	sumAddr := os.Getenv("SUMMARYADDR")
	sumAddrs := strings.Split(sumAddr, ",")
	sumAddrURLs := []*url.URL{}
	for i, _ := range sumAddrs {
		sumAddrURL, err := url.Parse(sumAddrs[i])
		if err != nil {
			fmt.Printf("error parsing summary URLs: %v\n", err)
			os.Exit(1)
		}
		sumAddrURLs = append(sumAddrURLs, sumAddrURL)
	}
	summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(sumAddrURLs)}

	// Receive address(es) for messages microservice(s) and insert into CustomDirector
	messagesAddr := os.Getenv("MESSAGESADDR")
	messagesAddrs := strings.Split(messagesAddr, ",")
	messagesAddrURLs := []*url.URL{}
	for i, _ := range messagesAddrs {
		messagesAddrURL, err := url.Parse(messagesAddrs[i])
		if err != nil {
			fmt.Printf("error parsing messages URLs: %v\n", err)
			os.Exit(1)
		}
		messagesAddrURLs = append(messagesAddrURLs, messagesAddrURL)
	}
	messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(messagesAddrURLs)}

	/*- Create a new mux for the web server.*/
	mux := http.NewServeMux()
	/*- Tell the mux to call your handlers.SummaryHandler function
	when the "/v1/summary" URL path is requested.*/

	mux.Handle("/v1/summary", summaryProxy)
	mux.Handle("/v1/channels", messagesProxy)
	mux.Handle("/v1/channels/", messagesProxy)
	mux.Handle("/v1/messages/", messagesProxy)
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler)
	mux.HandleFunc("/v1/ws", ctx.WebsocketConnectionHandler)
	/*- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/
	wrappedMux := handlers.NewEnsureCORS(handlers.NewEnsureAuth(mux, signingKey, redisStore))
	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))

}
