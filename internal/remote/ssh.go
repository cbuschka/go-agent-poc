package remote

import (
	"fmt"
	"github.com/cbuschka/go-agent-poc/internal/local"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"strings"
)

type RemoteSession struct {
	sshConn *ssh.Client
}

func OpenSshSession(username string, hostname string, port int) (Session, error) {
	session := RemoteSession{}
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			getSSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", hostname, port)
	sshConn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	session.sshConn = sshConn
	return &session, nil
}

func (s *RemoteSession) PlaceAgent(filepath string, agentBytes []byte) error {

	sftpSession, err := sftp.NewClient(s.sshConn)
	if err != nil {
		return fmt.Errorf("failed to start sftp session: %v", err)
	}
	defer sftpSession.Close()

	file, err := sftpSession.Create(filepath)
	if err != nil {
		return err
	}

	_, err = file.Write(agentBytes)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = sftpSession.Chmod(filepath, 0700)
	if err != nil {
		return err
	}

	return nil
}

func (s *RemoteSession) Close() error {
	if s.sshConn != nil {
		return s.sshConn.Close()
	}

	return nil
}

func getSSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

type RemoteAgent struct {
	sshSession *ssh.Session
	conn       *local.LocalConnection
}

func (a *RemoteAgent) GetConnection() (net.Conn, error) {
	return net.Conn(a.conn), nil
}

func (a *RemoteAgent) Close() error {
	if a.sshSession != nil {
		return a.sshSession.Close()
	}

	return nil
}

func (s *RemoteSession) StartAgent(command ...string) (Agent, error) {

	sshSession, err := s.sshConn.NewSession()
	if err != nil {
		return nil, err
	}

	sshSession.Stderr = NewPrefixedLineWriter("agent:")

	stdinWr, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutRd, err := sshSession.StdoutPipe()
	if err != nil {
		return nil, err
	}

	commandLine := strings.Join(command, " ")
	err = sshSession.Start(commandLine)
	if err != nil {
		return nil, err
	}

	remoteAgent := RemoteAgent{sshSession: sshSession}
	remoteAgent.conn = local.NewConnection(stdoutRd, stdinWr, nil)
	return &remoteAgent, nil
}
