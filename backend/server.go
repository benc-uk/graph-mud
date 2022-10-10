package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Simple server that listens on port 8080 and returns a JSON response
// with the current time.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hello, World!",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
