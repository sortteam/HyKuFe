package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

type SshClient struct {
	Connection *ssh.Client
	Session *ssh.Session
}

func (sc *SshClient) NewSshClient(address string, port int64, password string) error {
	sshConfig := &ssh.ClientConfig{
		Config:            ssh.Config{},
		User:              "sort-server-1",
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:           time.Second * 10,
	}

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, port), sshConfig)
	if err != nil {
		return fmt.Errorf("Failed to dial : %s", err)
	}
	sc.Connection = connection

	session, err := connection.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}
	sc.Session = session

	modes := ssh.TerminalModes{
		ssh.ECHO:		0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdout for session : %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session : %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	return nil
}

func (sc *SshClient) CommandExecution(command string) error {
	if sc.Session == nil || sc.Connection == nil {
		return fmt.Errorf("No Connection Or Session.")
	}

	if err := sc.Session.Run(command); err != nil {
		return fmt.Errorf("Command Execution Error : %v", err)
	}

	return nil
}