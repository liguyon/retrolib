package login

type ServerState int

const (
	ServerOffline ServerState = iota
	ServerOnline
	ServerStarting
)

type ServerType int

const (
	ServerClassic ServerType = iota
	ServerHardcore
	ServerEvent
)

type Server struct {
	ID         int
	State      ServerState
	Type       ServerType
	Completion int // TODO: what is it? n chars online?
	CanLogIn   bool
}

type ServerWithCharacters struct {
	ServerID   int
	NCharacter int
}
