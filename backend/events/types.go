package events

import (
	"nano-realms/backend/graph"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Processor struct {
	db    neo4j.Driver
	graph *graph.GraphService
}

type Processable interface {
	PreHook(*Processor) error
	Process(*Processor) error
	PostHook(*Processor) error
}

type BaseEvent struct {
	PreHookFunc  func() error
	PostHookFunc func() error
}

type CreateEvent struct {
	Type  NodeType
	Props map[string]any
	BaseEvent
}

type MoveEvent struct {
	NodeType  NodeType
	NodeProp  string
	NodeValue any

	DestType  NodeType
	DestProp  string
	DestValue any

	Relation Relationship
	BaseEvent
}

type DestroyEvent struct {
	NodeType NodeType
	Prop     string
	Value    any
	BaseEvent
}

type NodeType string

const (
	TypePlayer   NodeType = "Player"
	TypeLocation NodeType = "Location"
	TypeItem     NodeType = "Item"
	TypeSystem   NodeType = "System"
)

type Relationship string

const (
	RelIn     Relationship = "IN"
	RelHolds  Relationship = "HOLDS"
	RelEquips Relationship = "EQUIPS"
)
