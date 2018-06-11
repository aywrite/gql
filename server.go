// Command graphql-go-example starts an HTTP GraphQL API server which is backed by data
// against the https://swapi.co REST API.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aywrite/gql/handler"
	"github.com/aywrite/gql/loader"
	"github.com/aywrite/gql/models"
	"github.com/aywrite/gql/resolver"
	"github.com/aywrite/gql/schema"
	graphql "github.com/graph-gophers/graphql-go"
)

func main() {
	// Tweak configuration values here.
	var (
		addr              = ":8000"
		readHeaderTimeout = 1 * time.Second
		writeTimeout      = 10 * time.Second
		idleTimeout       = 90 * time.Second
		maxHeaderBytes    = http.DefaultMaxHeaderBytes
	)

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	root := &resolver.Resolver{}

	db := models.DB{
		ConnStr: "user=world dbname=world-db password=world123 sslmode=disable port=5433",
	}

	// Create the request handler; inject dependencies.
	h := handler.GraphQL{
		// Parse and validate schema. Panic if unable to do so.
		Schema:  graphql.MustParseSchema(schema.String(), root),
		Loaders: loader.Initialize(db),
	}

	// Register handlers to routes.
	mux := http.NewServeMux()
	mux.Handle("/", handler.GraphiQL{})
	mux.Handle("/graphql/", h)
	mux.Handle("/graphql", h) // Register without a trailing slash to avoid redirect.

	// Configure the HTTP server.
	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	// Begin listeing for requests.
	log.Printf("Listening for requests on %s", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		log.Println("server.ListenAndServe:", err)
	}

	// TODO: intercept shutdown signals for cleanup of connections.
	log.Println("Shut down.")
}
