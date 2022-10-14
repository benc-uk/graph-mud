package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"nano-realms-backend/api"
	"nano-realms-backend/messaging"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", api.HealthEndpoint).Methods("GET")
	router.HandleFunc("/", api.HealthEndpoint).Methods("GET")

	router.HandleFunc("/connect", messaging.UserConnect)

	router.HandleFunc("/player/{id}", api.GetPlayer).Methods("GET")

	// Bypass CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      cors.Handler(router),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go messaging.SenderTestLoop()

	log.Println("### Graph MUD - backend API, listening on port:", port)
	log.Fatal(srv.ListenAndServe())
}
