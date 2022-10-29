// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"nano-realms/backend/events"
	"nano-realms/backend/game"
	"nano-realms/pkg/auth"
	"nano-realms/pkg/problem"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// All routes we need should be registered here
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/player", auth.JWTValidator(api.getPlayer)).Methods("GET")
	router.HandleFunc("/player", auth.JWTValidator(api.newPlayer)).Methods("POST")
	router.HandleFunc("/player", auth.JWTValidator(api.deletePlayer)).Methods("DELETE")
	router.HandleFunc("/cmd", auth.JWTValidator(api.playerCommand)).Methods("POST")
	router.HandleFunc("/player/location", auth.JWTValidator(api.playerLocation)).Methods("GET")
}

// Get existing player details
func (api API) getPlayer(resp http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("X-Username")

	res, err := api.graph.QuerySingleNode("MATCH (p:Player {username: $p0}) RETURN p", []string{username})
	if err != nil {
		problem.Wrap(500, req.RequestURI, "username", err).Send(resp)
		return
	}
	if res == nil {
		problem.Wrap(404, req.RequestURI, "username", errors.New("Player not found")).Send(resp)
		return
	}

	player := &game.Player{}
	_ = mapstructure.Decode(res.Props, player)

	resp.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(player)
	_, _ = resp.Write(json)
}

// Handle new player creation
func (api API) newPlayer(resp http.ResponseWriter, req *http.Request) {
	var newPlayer game.NewPlayer
	username := req.Header.Get("X-Username")
	if username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username header")).Send(resp)
		return
	}

	// check player exists
	exists, err := api.graph.NodeExists("Player", "username", username)
	if err != nil || exists {
		problem.Wrap(409, req.RequestURI, "username", errors.New("Player already exists")).Send(resp)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&newPlayer)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "new", err).Send(resp)
		return
	}
	newPlayer.Username = username

	err = api.event.Process(events.CreateEvent{
		Type: events.TypePlayer,
		Props: map[string]interface{}{
			"username":    newPlayer.Username,
			"name":        newPlayer.Name,
			"class":       newPlayer.Class,
			"description": newPlayer.Description,
		},
	})

	if err != nil {
		problem.Wrap(500, req.RequestURI, "new", err).Send(resp)
		return
	}

	err = api.event.Process(events.MoveEvent{
		NodeType:  events.TypePlayer,
		NodeProp:  "username",
		NodeValue: newPlayer.Username,
		DestType:  events.TypeLocation,
		DestProp:  "name",
		DestValue: "lobby",
		Relation:  "IN",
	})
	if err != nil {
		problem.Wrap(500, req.RequestURI, "new", err).Send(resp)
		return
	}

	resp.WriteHeader(http.StatusCreated)
}

// Handle player removal
func (api API) deletePlayer(resp http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("X-Username")
	if username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username header")).Send(resp)
		return
	}

	err := api.event.Process(events.DestroyEvent{
		NodeType: events.TypePlayer,
		Prop:     "username",
		Value:    username,
	})

	if err != nil {
		problem.Wrap(500, req.RequestURI, "delete", err).Send(resp)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}

func (api API) playerCommand(resp http.ResponseWriter, req *http.Request) {
	var cmd game.Command
	username := req.Header.Get("X-Username")
	if username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username header")).Send(resp)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&cmd)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "command", err).Send(resp)
		return
	}

	log.Printf("### User %s: %s", username, cmd.Text)

	err = api.command.Handle(username, cmd.Text)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "command", err).Send(resp)
		return
	}

	_, _ = resp.Write([]byte("OK"))
}

func (api API) playerLocation(resp http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("X-Username")
	if username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username header")).Send(resp)
		return
	}
	res, err := api.graph.QuerySingleNode("MATCH (p:Player {username:$p0})-[IN]->(l:Location) RETURN l", []string{username})
	if err != nil {
		problem.Wrap(500, req.RequestURI, "username", err).Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(res.Props)
	_, _ = resp.Write(json)
}
