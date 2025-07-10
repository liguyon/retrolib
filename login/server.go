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

// Server represents a Retro game server or host (e.g. Jiva, Eratz, ...).
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
