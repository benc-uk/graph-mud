package graph

import (
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var Service *GraphService

type GraphService struct {
	db neo4j.Driver
}

func InitService(db neo4j.Driver) {
	Service = &GraphService{
		db: db,
	}
}

func (s *GraphService) NodeExists(label string, propName string, probValue string) (bool, error) {
	query := fmt.Sprintf("MATCH (n:%s {%s: $p0}) RETURN n", label, propName)
	node, err := s.QuerySingleNode(query, []string{probValue})
	if err != nil {
		return false, err
	}

	return node != nil, nil
}

func (s *GraphService) QuerySingleNode(query string, paramArray []string) (*neo4j.Node, error) {
	params := make(map[string]interface{})
	for i, p := range paramArray {
		params[fmt.Sprintf("p%d", i)] = p
	}

	sess := s.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	val, err := sess.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		res, err := tx.Run(query, params)

		if err != nil {
			return nil, err
		}

		if res.Next() {
			return res.Record().Values[0], nil
		} else {
			return nil, nil
		}
	})

	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	node := val.(neo4j.Node)
	return &node, err
}

func (s *GraphService) QueryMultiNode(query string, paramArray []string) ([]neo4j.Node, error) {
	params := make(map[string]interface{})
	for i, p := range paramArray {
		params[fmt.Sprintf("p%d", i)] = p
	}

	sess := s.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	val, err := sess.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		res, err := tx.Run(query, params)

		if err != nil {
			return nil, err
		}

		var nodes []neo4j.Node

		for res.Next() {
			node := res.Record().Values[0].(neo4j.Node)
			nodes = append(nodes, node)
		}
		return nodes, nil
	})

	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	out := val.([]neo4j.Node)
	return out, err
}

func (s *GraphService) QueryMultiRelationship(query string, paramArray []string) ([]neo4j.Relationship, error) {
	params := make(map[string]interface{})
	for i, p := range paramArray {
		params[fmt.Sprintf("p%d", i)] = p
	}

	sess := s.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	val, err := sess.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		res, err := tx.Run(query, params)

		if err != nil {
			return nil, err
		}

		var rels []neo4j.Relationship

		for res.Next() {
			rel := res.Record().Values[0].(neo4j.Relationship)
			rels = append(rels, rel)
		}
		return rels, nil
	})

	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	out := val.([]neo4j.Relationship)
	return out, err
}

func (s *GraphService) GetSingleNodeById(id int64) (*neo4j.Node, error) {
	query := fmt.Sprintf("MATCH (n) WHERE ID(n) = %d RETURN n", id)
	return s.QuerySingleNode(query, []string{})
}

func (s *GraphService) GetPlayerLocation(username string) (*neo4j.Node, error) {
	query := "MATCH (p:Player {username: $p0})-[:IN]->(l:Location) RETURN l"
	loc, err := s.QuerySingleNode(query, []string{username})
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, errors.New("player location not found")
	}
	return loc, nil
}

func (s *GraphService) GetPlayersInLocation(locationName string) ([]neo4j.Node, error) {
	query := "MATCH (p:Player)-[:IN]->(l:Location {name: $p0}) RETURN p"
	return s.QueryMultiNode(query, []string{locationName})
}

func (s *GraphService) GetPlayer(username string) (*neo4j.Node, error) {
	query := "MATCH (p:Player {username: $p0}) RETURN p"
	player, err := s.QuerySingleNode(query, []string{username})
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, fmt.Errorf("player %s not found", username)
	}
	return player, nil
}
