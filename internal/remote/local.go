package remote

import (
	"github.com/cbuschka/go-agent-poc/internal/local"
	"io/ioutil"
	"net"
	"os/exec"
)

type LocalSession struct {
}

func OpenLocalSession() (Session, error) {
	return &LocalSession{}, nil
}

func (s *LocalSession) PlaceAgent(filepath string, agentBytes []byte) error {
	return ioutil.WriteFile(filepath, agentBytes, 0700)
}

func (s *LocalSession) Close() error {
	return nil
}

type LocalAgent struct {
	conn *local.LocalConnection
	cmd  *exec.Cmd
}

func (a *LocalAgent) GetConnection() (net.Conn, error) {
	return net.Conn(a.conn), nil
}

func (a *LocalAgent) Close() error {
	if a.cmd != nil {
		if a.cmd.ProcessState != nil && !a.cmd.ProcessState.Exited() {
			_ = a.cmd.Process.Kill()
		}
	}
	return nil
}

func (s *LocalSession) StartAgent(command ...string) (Agent, error) {

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = NewPrefixedLineWriter("agent:")

	stdinWr, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutRd, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	remoteAgent := LocalAgent{cmd: cmd}
	remoteAgent.conn = local.NewConnection(stdoutRd, stdinWr, nil)
	return &remoteAgent, nil
}
