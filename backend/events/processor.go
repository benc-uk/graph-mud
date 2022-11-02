package events

import (
	"fmt"
	"nano-realms/backend/graph"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func NewProcessor(dbDriver neo4j.Driver, graph *graph.GraphService) *Processor {
	return &Processor{
		db:    dbDriver,
		graph: graph,
	}
}

func (p *Processor) Process(event Processable) error {
	err := event.PreHook(p)
	if err != nil {
		return err
	}
	err = event.Process(p)
	if err != nil {
		return err
	}
	err = event.PostHook(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) createNode(nodeType string, props any) error {
	sess := p.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(`
		  CREATE (p:%s $props)`, nodeType),
			map[string]any{"props": props})

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func (p *Processor) moveNode(nodeType string, nodeProp string, nodeValue any, rel string, destType string, destProp string, destValue any) error {
	sess := p.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := fmt.Sprintf(` 
			MATCH (n:%s {%s:$nv})
			OPTIONAL MATCH (n)-[r:%s]->(s) 
			MATCH (d:%s {%s:$dv}) 
			DELETE r
			CREATE (n)-[:%s]->(d)`,
			nodeType, nodeProp, rel, destType, destProp, rel)
		params := map[string]any{"nv": nodeValue, "dv": destValue}
		_, err := tx.Run(query, params)

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

func (p *Processor) deleteNode(nodeType, nodeProp string, nodeValue any) error {
	sess := p.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		query := fmt.Sprintf("MATCH (n:%s {%s: $v}) DETACH DELETE n", nodeType, nodeProp)
		_, err := tx.Run(query, map[string]any{"v": nodeValue})

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}
