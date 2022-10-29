package commands

import (
	"errors"
	"fmt"
	"nano-realms/backend/events"
	"nano-realms/backend/graph"
	"nano-realms/backend/messaging"
	"strings"

	"golang.org/x/exp/slices"
)

type Handler struct {
	graph *graph.GraphService
	event *events.Processor
}

func NewHandler(graph *graph.GraphService, processor *events.Processor) *Handler {
	return &Handler{
		graph: graph,
		event: processor,
	}
}

func (h *Handler) Handle(username string, cmd string) error {
	c := strings.ToLower(cmd)
	c = strings.TrimSpace(c)
	cParts := strings.Split(c, " ")

	if slices.Contains([]string{"look", "where", "l"}, cParts[0]) {
		res, err := h.graph.QuerySingleNode("MATCH (p:Player {username:$p0})-[IN]->(l:Location) RETURN l", []string{username})
		if err != nil {
			return err
		}
		locDesc := res.Props["description"].(string)

		messaging.SendToUser(username, fmt.Sprintf("You can see: %s", locDesc), "command", "look")
		return nil
	}

	if slices.Contains([]string{"north", "south", "east", "west", "n", "s", "e", "w"}, cParts[0]) {
		exits, err := h.graph.QueryMultiRelationship(`MATCH (:Player {username:$p0})-[:IN]->(l:Location) MATCH (l)-[r]->(:Location) RETURN r`, []string{username})
		if err != nil || exits == nil {
			return errors.New("Exits not found")
		}

		for _, exit := range exits {
			// Find the exit that matches the direction
			if strings.ToLower(exit.Type)[0:1] == cParts[0][0:1] {
				loc, err := h.graph.GetSingleNodeById(exit.EndId)
				if err != nil || loc == nil {
					return errors.New("Location not found")
				}

				return h.event.Process(events.MovePlayerEvent(username, loc.Props["name"].(string)))
			}
		}

		messaging.SendToUser(username, "You can't go that way", "command", "blocked")
		return nil
	}

	if cParts[0] == "$play" {
		e := events.MovePlayerEvent(username, "gameEntry")
		e.DestProp = "gameEntry"
		e.DestValue = true
		return h.event.Process(e)
	}

	if cParts[0] == "$lobby" {
		return h.event.Process(events.MovePlayerEvent(username, "lobby"))
	}

	messaging.SendToUser(username, "Invalid command: "+cmd, "command", "invalid")
	return nil
}
