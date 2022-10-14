package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte(`{"alive": true}`))
	})

	router.HandleFunc("/players", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte(`{"p1": "hello"}`))
	})

	router.HandleFunc("/connect", wsConnect)

	// add cors
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

	go sender()

	log.Println("### Graph MUD - backend API, listening on port:", port)
	log.Fatal(srv.ListenAndServe())
}
