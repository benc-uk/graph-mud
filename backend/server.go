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

	"github.com/benc-uk/go-rest-api/pkg/api"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/env"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	_ "github.com/joho/godotenv/autoload"
)

// API type is a wrap of the common base API with local implementation
type API struct {
	*api.Base
	event   *events.Processor
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

	router := chi.NewRouter()

	graph.InitService(dbDriver)
	eventProcessor := events.NewProcessor(dbDriver)
	api := API{
		api.NewBase(serviceName, version, buildInfo, healthy),
		eventProcessor,
		commands.NewHandler(eventProcessor),
	}

	// Some basic middleware
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// Some custom middleware for CORS & JWT username
	router.Use(api.SimpleCORSMiddleware)

	router.Group(func(appRouter chi.Router) {
		clientID := os.Getenv("AUTH_CLIENT_ID")
		if clientID == "" {
			log.Println("### 🚨 No AUTH_CLIENT_ID set, skipping auth validation")
		} else {
			log.Println("### 🔐 Auth enabled, validating JWT tokens")
			jwtValidator := auth.NewJWTValidator(clientID,
				"https://login.microsoftonline.com/common/discovery/v2.0/keys",
				"Play.Game")

			appRouter.Use(jwtValidator.Middleware)
		}

		appRouter.Use(api.JWTRequestEnricher("username", "preferred_username"))
		api.addRoutes(appRouter)
	})

	router.Group(func(publicRouter chi.Router) {
		// Add Prometheus metrics endpoint, must be before the other routes
		api.AddMetricsEndpoint(publicRouter, "metrics")

		// Add optional root, health & status endpoints
		api.AddHealthEndpoint(publicRouter, "health")
		api.AddStatusEndpoint(publicRouter, "status")
		api.AddOKEndpoint(publicRouter, "")
	})

	// For websocket connections & messaging
	messaging.Version = version
	messaging.AddRoutes(router)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	worldOK, err := graph.Service.NodeExists("Location", "name", "lobby")
	if err != nil {
		log.Fatal(err)
	}
	if !worldOK {
		api.Healthy = false
		log.Println("### 💥 SEVERE! Database not configured, please run load the realm data!")
	} else {
		log.Println("### 🌍 Realm data appears OK")
	}

	log.Printf("### 🌐 Nano Realms Backend API, listening on port: %d", serverPort)
	log.Printf("### 🚀 Build details: v%s (%s)", version, buildInfo)
	log.Fatal(srv.ListenAndServe())
}
