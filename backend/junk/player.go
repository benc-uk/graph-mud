package game

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type Player struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Class       string `json:"class"`
	Description string `json:"description"`
}

type PlayerService struct {
	db neo4j.Driver
}

func NewPlayerService(dbDriver neo4j.Driver) *PlayerService {
	return &PlayerService{
		db: dbDriver,
	}
}

func (s *PlayerService) Get(username string) (*Player, error) {
	sess := s.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	p, err := sess.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		res, err := tx.Run(`
		  MATCH (p:Player {username: $u}) 
		  RETURN p`,
			map[string]any{"u": username})

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

	if p == nil {
		return nil, errors.New("Player not found")
	}

	node := p.(dbtype.Node)
	return &Player{
		Username:    node.Props["username"].(string),
		Name:        node.Props["name"].(string),
		Class:       node.Props["class"].(string),
		Description: node.Props["description"].(string),
	}, nil
}

// Create a new player, and place them in the lobby area
// Username must be unique
func (s *PlayerService) Create(player NewPlayer) error {
	sess := s.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		res, err := tx.Run(`
				MATCH (lobby:Location {name: 'lobby'})
				MERGE (p:Player {username: $u})
				MERGE (p)-[:IN]->(lobby)
				ON CREATE
					SET p.description = $d, p.class = $c, p.name = $n
				RETURN p;`,
			map[string]any{"u": player.Username, "n": player.Name, "c": player.Class, "d": player.Description})

		if err != nil {
			return nil, err
		}

		if res.Next() {
			return res.Record().Values[0], nil
		} else {
			return nil, nil
		}
	})

	return err
}

// Delete a player
func (s *PlayerService) Delete(username string) error {
	sess := s.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err := sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
				MATCH (p:Player {username: $u})
				DETACH DELETE p`,
			map[string]any{"u": username})

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}
