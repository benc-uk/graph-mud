package game

type NewPlayer struct {
	Username    string
	Name        string
	Class       string
	Description string
}

type Player struct {
	Username    string `mapstructure:"username"`
	Name        string `mapstructure:"name"`
	Class       string `mapstructure:"class"`
	Description string `mapstructure:"description"`
}
