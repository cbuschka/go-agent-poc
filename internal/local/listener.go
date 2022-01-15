package local

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type Listener struct {
	closed   bool
	lock     *sync.Mutex
	currConn *LocalConnection
}

func NewListener() *Listener {
	lock := sync.Mutex{}
	lis := Listener{closed: false, lock: &lock}
	return &lis
}

func (lis *Listener) Accept() (net.Conn, error) {
	if lis.closed {
		return nil, fmt.Errorf("listener closed")
	}

	lis.lock.Lock()

	if lis.closed {
		lis.lock.Unlock()
		return nil, fmt.Errorf("listener closed")
	}

	newConn := NewConnection(os.Stdin, os.Stdout, lis)
	lis.currConn = newConn
	return newConn, nil
}

func (lis *Listener) Close() error {
	lis.closed = true
	return nil
}

func (lis *Listener) Addr() net.Addr {
	return NewAddr("local")
}

func (lis *Listener) connClosed(c *LocalConnection) {
	if lis.currConn == c {
		lis.currConn = nil
		lis.lock.Unlock()
	}
}
