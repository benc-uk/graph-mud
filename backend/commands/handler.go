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
	event *events.Processor
}

func NewHandler(processor *events.Processor) *Handler {
	return &Handler{
		event: processor,
	}
}

func (h *Handler) Handle(username string, cmd string) error {
	c := strings.ToLower(cmd)
	c = strings.TrimSpace(c)
	cParts := strings.Split(c, " ")

	if slices.Contains([]string{"look", "where", "l"}, cParts[0]) {
		res, err := graph.Service.GetPlayerLocation(username)
		if err != nil {
			return err
		}
		locDesc := res.Props["description"].(string)

		messaging.SendToUser(username, fmt.Sprintf("You can see: %s", locDesc), "command", "look")
		return nil
	}

	if slices.Contains([]string{"north", "south", "east", "west", "n", "s", "e", "w"}, cParts[0]) {
		exits, err := graph.Service.QueryMultiRelationship(`MATCH (:Player {username:$p0})-[:IN]->(l:Location) MATCH (l)-[r]->(:Location) RETURN r`, []string{username})
		if err != nil || exits == nil {
			return errors.New("Exits not found")
		}

		for _, exit := range exits {
			// Find the exit that matches the direction
			if strings.ToLower(exit.Type)[0:1] == cParts[0][0:1] {
				loc, err := graph.Service.GetSingleNodeById(exit.EndId)
				if err != nil || loc == nil {
					return errors.New("Location not found")
				}

				return h.event.Process(events.NewPlayerMoveEvent(username, loc.Props["name"].(string)))
			}
		}

		messaging.SendToUser(username, "You can't go that way", "command", "blocked")
		return nil
	}

	if cParts[0] == "$play" {
		e := events.NewPlayerMoveEvent(username, "gameEntry")
		e.DestProp = "gameEntry"
		e.DestValue = true
		return h.event.Process(e)
	}

	if cParts[0] == "$lobby" {
		return h.event.Process(events.NewPlayerMoveEvent(username, "lobby"))
	}

	if cParts[0] == "say" || cParts[0] == "speak" {
		if len(cParts) < 2 {
			messaging.SendToUser(username, "Say what?", "command", "invalid")
			return nil
		}

		// Get the player's location
		location, err := graph.Service.GetPlayerLocation(username)
		if err != nil {
			return err
		}
		player, err := graph.Service.GetPlayer(username)
		if err != nil {
			return err
		}

		// Send message to other players in the location
		text := fmt.Sprintf("%s says \"%s\"", player.Props["name"], strings.Join(cParts[1:], " "))
		msg := messaging.NewGameMessage(text, "command", "say")
		messaging.SendToAllUsersInLocation(location.Props["name"].(string), msg, username)
		return nil
	}

	messaging.SendToUser(username, "Invalid command: "+cmd, "command", "invalid")
	return nil
}
