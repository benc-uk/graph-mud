package events

import (
	"fmt"

	"nano-realms/backend/messaging"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Processor struct {
	db neo4j.Driver
}

func NewProcessor(dbDriver neo4j.Driver) *Processor {
	return &Processor{
		db: dbDriver,
	}
}

func (p *Processor) Process(event any) error {
	// check  event is CreateEvent
	ce, ok := event.(CreateEvent)
	if ok {
		return p.createNode(ce)
	}
	me, ok := event.(MoveEvent)
	if ok {
		messaging.SendToUser("oooooooooooooo", "You have moved to the"+me.DestValue, "system", "system")
		return p.moveNode(me)
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

func (event *CreateEvent) Notify() error {
	return nil
}
