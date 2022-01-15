package local

// from https://github.com/bithavoc/hellogrpcstdin/blob/master/common/joint.go

import (
	"encoding/hex"
	"io"
	"net"
	"time"
)

type Addr struct {
	s string
}

func NewAddr(s string) *Addr {
	return &Addr{s}
}
func (a *Addr) Network() string {
	return "stdio"
}

func (a *Addr) String() string {
	return a.s
}

type LocalConnection struct {
	in       io.Reader
	out      io.Writer
	closed   bool
	local    *Addr
	remote   *Addr
	listener *Listener
}

func NewConnection(in io.Reader, out io.Writer, listener *Listener) *LocalConnection {
	return &LocalConnection{
		local:    NewAddr("local"),
		remote:   NewAddr("remote"),
		in:       hex.NewDecoder(in),
		out:      hex.NewEncoder(out),
		listener: listener,
	}
}

func (s *LocalConnection) LocalAddr() net.Addr {
	return s.local
}

func (s *LocalConnection) RemoteAddr() net.Addr {
	return s.remote
}

func (s *LocalConnection) Read(b []byte) (n int, err error) {
	return s.in.Read(b)
}

func (s *LocalConnection) Write(b []byte) (n int, err error) {
	return s.out.Write(b)
}

func (s *LocalConnection) Close() error {
	s.closed = true
	if s.listener != nil {
		s.listener.connClosed(s)
	}
	return nil
}

func (s *LocalConnection) SetDeadline(t time.Time) error {
	return nil
}

func (s *LocalConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (s *LocalConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
