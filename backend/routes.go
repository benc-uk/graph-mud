// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"net/http"

	"nano-realms/backend/game"
	"nano-realms/pkg/auth"
	"nano-realms/pkg/problem"

	"github.com/gorilla/mux"
)

// All routes we need should be registered here
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/player", auth.JWTValidator(api.getPlayer)).Methods("GET")
	router.HandleFunc("/player", auth.JWTValidator(api.newPlayer)).Methods("POST")
}

// ============================================================================

func (api API) getPlayer(resp http.ResponseWriter, req *http.Request) {
	username := "ben"
	p, err := api.player.Get(username)
	if err != nil {
		problem.New500(req.RequestURI, username, err).Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(p)
	_, _ = resp.Write(json)
}

// ============================================================================

func (api API) newPlayer(resp http.ResponseWriter, req *http.Request) {
	var newPlayer game.NewPlayer
	err := json.NewDecoder(req.Body).Decode(&newPlayer)
	if err != nil {
		problem.New500(req.RequestURI, "NewPlayer", err).Send(resp)
		return
	}

	err = api.player.Create(newPlayer)
	if err != nil {
		problem.New500(req.RequestURI, "NewPlayer", err).Send(resp)
		return
	}

	resp.WriteHeader(http.StatusCreated)
}
