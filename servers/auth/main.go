package main

import (
	"GroupPoop/servers/auth/handlers"
	"GroupPoop/servers/auth/models/users"
	"GroupPoop/servers/auth/sessions"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

//main is the main entry point for the server
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

	// connect to mysql database
	// Bradley: set database password as an environment variable
	// and connect to MySQL database
	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Printf("error pinging database: %v\n", err)
	} else {
		log.Printf("successfully connected!\n")
	}

	userStore := users.NewMySQLStore(db)
	signingKey := os.Getenv("SESSIONKEY")

	ctx := handlers.NewHandlerContext(signingKey, redisStore, userStore)

	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80" */
	addr := os.Getenv("AUTHPORT")
	if len(addr) == 0 {
		addr = ":443"
	}

	/*- Create a new mux for the web server.*/
	mux := http.NewServeMux()
	/*- Tell the mux to call your handlers.SummaryHandler function
	when the "/v1/summary" URL path is requested.*/

	mux.HandleFunc("/login", ctx.SessionsHandler)
	mux.HandleFunc("/login/", ctx.SpecificSessionHandler)
	/*- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/
	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))

}
