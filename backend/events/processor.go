package events

import (
	"fmt"
	"nano-realms/backend/graph"
	"nano-realms/backend/messaging"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Processor struct {
	db    neo4j.Driver
	graph *graph.GraphService
}

func NewProcessor(dbDriver neo4j.Driver, graph *graph.GraphService) *Processor {
	return &Processor{
		db:    dbDriver,
		graph: graph,
	}
}

func (p *Processor) Process(event any) error {
	ce, ok := event.(CreateEvent)
	if ok {
		return p.createNode(ce)
	}

	me, ok := event.(MoveEvent)
	if ok {

		err := p.moveNode(me)
		if err != nil {
			return err
		}

		if me.NodeType == TypePlayer {
			username := me.NodeValue.(string)
			res, err := p.graph.QuerySingleNode("MATCH (p:Player {username:$p0})-[IN]->(l:Location) RETURN l", []string{username})
			if err != nil {
				return err
			}
			locDesc := res.Props["description"].(string)

			messaging.SendToUser(username, "You move into: "+locDesc, "server", "move")
		}
	}

	de, ok := event.(DestroyEvent)
	if ok {
		return p.deleteNode(de)
	}

	return nil
}

func (p *Processor) createNode(event CreateEvent) error {
	sess := p.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(`
		  CREATE (p:%s $props)`, event.Type),
			map[string]any{"props": event.Props})

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func (p *Processor) moveNode(e MoveEvent) error {
	sess := p.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := fmt.Sprintf(` 
			MATCH (n:%s {%s:$nv})
			OPTIONAL MATCH (n)-[r:%s]->(s) 
			MATCH (d:%s {%s:$dv}) 
			DELETE r
			CREATE (n)-[:%s]->(d)`,
			e.NodeType, e.NodeProp, e.Relation, e.DestType, e.DestProp, e.Relation)
		params := map[string]any{"nv": e.NodeValue, "dv": e.DestValue}
		_, err := tx.Run(query, params)

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func (p *Processor) deleteNode(e DestroyEvent) error {
	sess := p.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := fmt.Sprintf("MATCH (n:%s {%s: $v}) DETACH DELETE n", e.NodeType, e.Prop)
		_, err := tx.Run(query, map[string]any{"v": e.Value})

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func MovePlayerEvent(username, dest string) MoveEvent {
	return MoveEvent{
		NodeType:  TypePlayer,
		NodeProp:  "username",
		NodeValue: username,
		DestType:  TypeLocation,
		DestProp:  "name",
		DestValue: dest,
		Relation:  RelIn,
	}
}
