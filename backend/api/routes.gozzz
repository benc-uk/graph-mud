package api

import (
	"encoding/json"
	"fmt"
	"nano-realms-backend/game"
	"net/http"
)

func HealthEndpoint(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	_, _ = resp.Write([]byte(`{"alive": true}`))
}

func GetPlayer(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	p := game.Player{
		Username: "ggg",
		Class:    "sffsfs",
	}
	sendJSON(resp, p)
}

func sendJSON(resp http.ResponseWriter, thing any) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(thing)
	//_, _ = resp.Write(jsonString)
}

func sendError(resp http.ResponseWriter, message string, status int) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)
	_, _ = resp.Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
}
