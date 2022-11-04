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
	// Update the graph database
	return p.moveNode(string(e.NodeType), e.NodeProp, e.NodeValue, string(e.Relation), string(e.DestType), e.DestProp, e.DestValue)
}

func NewPlayerMoveEvent(username string, dest string) *MoveEvent {
	e := &MoveEvent{
		NodeType:  TypePlayer,
		NodeProp:  "username",
		NodeValue: username,
		DestType:  TypeLocation,
		DestProp:  "name",
		DestValue: dest,
		Relation:  RelIn,
	}

	playerRes, err := graph.Service.GetPlayer(username)
	if err != nil {
		return nil
	}
	moverName := playerRes.Props["name"].(string)

	e.PreHookFunc = func() error {
		locationRes, err := graph.Service.GetPlayerLocation(username)
		if err != nil {

			return err
		}

		// Send message to other players on leaving, if you are in a location
		msg := messaging.NewGameMessage(moverName+" leaves", "command", "move")
		messaging.SendToAllUsersInLocation(locationRes.Props["name"].(string), msg, username)

		return nil
	}

	e.PostHookFunc = func() error {
		// Location will have changed
		locationRes, err := graph.Service.GetPlayerLocation(username)
		if err != nil {
			return err
		}

		// Send message to other players on entering, if you are in a location
		locDesc := locationRes.Props["description"].(string)
		messaging.SendToUser(username, "You move into: "+locDesc, "server", "move")

		msg := messaging.NewGameMessage(moverName+" enters", "command", "move")
		messaging.SendToAllUsersInLocation(locationRes.Props["name"].(string), msg, username)

		return nil
	}

	return e
}
