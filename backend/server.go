package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nano-realms/backend/commands"
	"nano-realms/backend/events"
	"nano-realms/backend/graph"
	"nano-realms/backend/messaging"
	"nano-realms/pkg/api"
	"nano-realms/pkg/auth"
	"nano-realms/pkg/env"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// API type is a wrap of the common base API with local implementation
type API struct {
	*api.Base
	event   *events.Processor
	graph   *graph.GraphService
	command *commands.Handler
}

var (
	healthy     = true               // Simple health flag
	version     = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName = "nano-realms"
	defaultPort = 8000
)

func main() {
	// Port to listen on, change the default as you see fit
	serverPort := env.GetEnvInt("PORT", defaultPort)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbUri := fmt.Sprintf("neo4j://%s:7687", dbHost)
	log.Printf("### 📥 Connecting to DB at %s", dbUri)
	var err error
	dbDriver, err := neo4j.NewDriver(dbUri, neo4j.NoAuth())
	if err != nil {
		log.Fatal(err)
	}
	err = dbDriver.VerifyConnectivity()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("### ✅ Connected to Neo4j!")
	}

	// Wrapper API with anonymous inner new Base API
	router := mux.NewRouter()

	graphService := graph.NewGraphService(dbDriver)
	eventProcessor := events.NewProcessor(dbDriver, graphService)
	api := API{
		api.NewBase(serviceName, version, buildInfo, healthy, router),
		eventProcessor,
		graphService,
		commands.NewHandler(graphService, eventProcessor),
	}

	// Main REST API routes
	api.addRoutes(router)

	// For websocket connections & messaging
	messaging.Version = version
	messaging.AddRoutes(router)

	// Extra routes for health and other features
	api.AddStatus(router)  // Add status and information endpoint
	api.AddLogging(router) // Add request logging
	api.AddHealth(router)  // Add health endpoint
	api.AddMetrics(router) // Expose metrics, in prometheus format
	api.AddRoot(router)    // Respond to root request with a simple 200 OK

	// Tell the auth middleware what scope to check when validating tokens
	auth.AppScopeName = "Play.Game"

	srv := &http.Server{
		Handler:      api.ConfigureCORSHandler(router),
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	log.Printf("### 🌐 Nano Realms Backend API, listening on port: %d", serverPort)
	log.Printf("### 🚀 Build details: v%s (%s)", version, buildInfo)
	log.Fatal(srv.ListenAndServe())
}
