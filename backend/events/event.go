package events

type CreateEvent struct {
	Type  NodeType
	Props map[string]any
}

type MoveEvent struct {
	NodeType  NodeType
	NodeProp  string
	NodeValue string

	DestType  NodeType
	DestProp  string
	DestValue string

	Relation Relationship
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
