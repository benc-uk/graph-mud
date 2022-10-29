package events

type CreateEvent struct {
	Type  NodeType
	Props map[string]any
}

type MoveEvent struct {
	NodeType  NodeType
	NodeProp  string
	NodeValue any

	DestType  NodeType
	DestProp  string
	DestValue any

	Relation Relationship
}

type DestroyEvent struct {
	NodeType NodeType
	Prop     string
	Value    any
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
