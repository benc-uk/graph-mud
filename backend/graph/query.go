package graph

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GraphService struct {
	db neo4j.Driver
}

func NewGraphService(db neo4j.Driver) *GraphService {
	return &GraphService{
		db: db,
	}
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
