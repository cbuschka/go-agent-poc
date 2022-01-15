package remote

import "net"

type Session interface {
	PlaceAgent(path string, agentBytes []byte) error

	StartAgent(command ...string) (Agent, error)

	Close() error
}

type Agent interface {
	GetConnection() (net.Conn, error)

	Close() error
}
