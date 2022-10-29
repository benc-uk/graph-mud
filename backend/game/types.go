package game

type NewPlayer struct {
	Username    string
	Name        string
	Class       string
	Description string
}

type Player struct {
	Username    string `mapstructure:"username" json:"username"`
	Name        string `mapstructure:"name" json:"name"`
	Class       string `mapstructure:"class" json:"class"`
	Description string `mapstructure:"description" json:"description"`
}

type Command struct {
	Text string
}
