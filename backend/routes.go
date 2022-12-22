// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// All application routes
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"nano-realms/backend/events"
	"nano-realms/backend/graph"

	"github.com/benc-uk/go-rest-api/pkg/problem"

	"github.com/go-chi/chi/v5"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type LocationResp struct {
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Exits       []string `json:"exits"`
}

type NewPlayer struct {
	Username    string
	Name        string
	Class       string
	Description string
	Gold        int
	Experience  int
	Rank        int
}

type PlayerResp struct {
	Username    string `mapstructure:"username" json:"username"`
	Name        string `mapstructure:"name" json:"name"`
	Class       string `mapstructure:"class" json:"class"`
	Description string `mapstructure:"description" json:"description"`
	Gold        int    `mapstructure:"gold" json:"gold"`
	Experience  int    `mapstructure:"experience" json:"experience"`
	Rank        int    `mapstructure:"rank" json:"rank"`
}

type Command struct {
	Text string
}

// All routes we need should be registered here
func (api API) addRoutes(router chi.Router) {
	router.Get("/player", api.getPlayer)
	router.Post("/player", api.newPlayer)
	router.Delete("/player", api.deletePlayer)
	router.Post("/cmd", api.executeCommand)
	router.Get("/player/location", api.playerLocation)
}

// Get existing player details
func (api API) getPlayer(resp http.ResponseWriter, req *http.Request) {
	username, ok := req.Context().Value("username").(string)
	if !ok || username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username")).Send(resp)
		return
	}

	res, err := graph.Service.QuerySingleNode("MATCH (p:Player {username: $p0}) RETURN p", []string{username})
	if err != nil {
		problem.Wrap(500, req.RequestURI, "username", err).Send(resp)
		return
	}
	if res == nil {
		problem.Wrap(404, req.RequestURI, "username", errors.New("Player not found")).Send(resp)
		return
	}

	player := &PlayerResp{}
	_ = mapstructure.Decode(res.Props, player)

	resp.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(player)
	_, _ = resp.Write(json)
}

// Handle new player creation
func (api API) newPlayer(resp http.ResponseWriter, req *http.Request) {
	var newPlayer NewPlayer
	username, ok := req.Context().Value("username").(string)

	if !ok || username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username")).Send(resp)
		return
	}

	// Check player exists
	exists, err := graph.Service.NodeExists("Player", "username", username)
	if err != nil || exists {
		problem.Wrap(409, req.RequestURI, "username", errors.New("Player already exists")).Send(resp)
		return
	}

	// Check for lobby as a way to check if world is setup
	lobbyExists, err := graph.Service.NodeExists("Location", "name", "lobby")
	if err != nil || !lobbyExists {
		problem.Wrap(500, req.RequestURI, "none", errors.New("Realm database has not been initialized! Contact the server admin")).Send(resp)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&newPlayer)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "new", err).Send(resp)
		return
	}
	newPlayer.Username = username

	err = api.event.Process(&events.CreateEvent{
		Type: events.TypePlayer,
		Props: map[string]interface{}{
			"username":    newPlayer.Username,
			"name":        newPlayer.Name,
			"class":       newPlayer.Class,
			"description": newPlayer.Description,
			"gold":        newPlayer.Gold,
			"experience":  newPlayer.Experience,
			"rank":        1,
		},
	})

	if err != nil {
		problem.Wrap(500, req.RequestURI, "new", err).Send(resp)
		return
	}

	err = api.event.Process(&events.MoveEvent{
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
	username, ok := req.Context().Value("username").(string)
	if !ok || username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username")).Send(resp)
		return
	}

	err := api.event.Process(&events.DestroyEvent{
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

func (api API) executeCommand(resp http.ResponseWriter, req *http.Request) {
	var cmd Command
	username, ok := req.Context().Value("username").(string)

	if !ok || username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username")).Send(resp)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&cmd)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "command", err).Send(resp)
		return
	}

	log.Printf("### User %s: %s", username, cmd.Text)

	// So this is pretty core bit of the code, where the REST API calls into the command handler
	err = api.command.Handle(username, cmd.Text)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "command", err).Send(resp)
		return
	}

	_, _ = resp.Write([]byte("OK"))
}

func (api API) playerLocation(resp http.ResponseWriter, req *http.Request) {
	username, ok := req.Context().Value("username").(string)
	if !ok || username == "" {
		problem.Wrap(400, req.RequestURI, "username", errors.New("Missing username header")).Send(resp)
		return
	}

	res, err := graph.Service.GetPlayerLocation(username)
	if err != nil {
		problem.Wrap(500, req.RequestURI, "username", err).Send(resp)
		return
	}

	exits, err := graph.Service.QueryMultiRelationship(`MATCH (:Player {username:$p0})-[:IN]->(l:Location) MATCH (l)-[r]->(:Location) RETURN r`, []string{username})
	if err != nil {
		problem.Wrap(500, req.RequestURI, "username", err).Send(resp)
		return
	}

	locationResp := &LocationResp{
		Description: res.Props["description"].(string),
		Name:        res.Props["name"].(string),
		Exits:       make([]string, 0),
	}

	caser := cases.Title(language.English)
	for _, e := range exits {
		locationResp.Exits = append(locationResp.Exits, caser.String(e.Type))
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(locationResp)
	_, _ = resp.Write(json)
}
