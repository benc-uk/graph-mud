package events

import (
	"nano-realms/backend/graph"
	"nano-realms/backend/messaging"
)

func (e *MoveEvent) PreHook(p *Processor) error {
	if e.PreHookFunc == nil {
		return nil
	}
	return e.PreHookFunc()
}

func (e *MoveEvent) PostHook(p *Processor) error {
	if e.PostHookFunc == nil {
		return nil
	}
	return e.PostHookFunc()
}

func (e *MoveEvent) Process(p *Processor) error {
	username := e.NodeValue.(string)
	playerRes, err := p.graph.GetPlayer(username)
	if err != nil {
		return err
	}
	moverName := playerRes.Props["name"].(string)

	// Send messages
	if e.NodeType == TypePlayer {
		locationRes, err := p.graph.GetPlayerLocation(username)
		if err != nil {
			return err
		}

		// Send message to other players on leaving, if you are in a location
		if locationRes != nil {
			players, err := p.graph.GetPlayersInLocation(locationRes.Props["name"].(string))
			if err != nil {
				return err
			}

			for _, player := range players {
				if u := player.Props["username"].(string); u != username {
					messaging.SendToUser(u, moverName+" leaves", "server", "player_enter")
				}
			}
		}
	}

	// Update the graph database
	err = p.moveNode(string(e.NodeType), e.NodeProp, e.NodeValue, string(e.Relation), string(e.DestType), e.DestProp, e.DestValue)
	if err != nil {
		return err
	}

	// Send messages
	if e.NodeType == TypePlayer {
		locationRes, err := p.graph.GetPlayerLocation(username)
		if err != nil {
			return err
		}

		if locationRes != nil {
			locDesc := locationRes.Props["description"].(string)
			messaging.SendToUser(username, "You move into: "+locDesc, "server", "move")

			// Send message to other players on entry to new location
			players, err := p.graph.GetPlayersInLocation(locationRes.Props["name"].(string))
			if err != nil {
				return err
			}

			for _, player := range players {
				if u := player.Props["username"].(string); u != username {
					messaging.SendToUser(u, moverName+" enters", "server", "player_enter")
				}
			}
		}
	}

	return nil
}

func NewPlayerMoveEvent(graph *graph.GraphService, username string, dest string) *MoveEvent {
	e := &MoveEvent{
		NodeType:  TypePlayer,
		NodeProp:  "username",
		NodeValue: username,
		DestType:  TypeLocation,
		DestProp:  "name",
		DestValue: dest,
		Relation:  RelIn,
	}

	e.PostHookFunc = func() error {
		username := e.NodeValue.(string)
		playerRes, err := graph.GetPlayer(username)
		if err != nil {
			return err
		}
		moverName := playerRes.Props["name"].(string)

		locationRes, err := graph.GetPlayerLocation(username)
		if err != nil {
			return err
		}

		// Send message to other players on leaving, if you are in a location
		if locationRes != nil {
			players, err := graph.GetPlayersInLocation(locationRes.Props["name"].(string))
			if err != nil {
				return err
			}

			for _, player := range players {
				if u := player.Props["username"].(string); u != username {
					messaging.SendToUser(u, moverName+" leaves", "server", "player_enter")
				}
			}
		}
		return nil
	}

	return e
}
