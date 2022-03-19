package sshx

import (
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Command struct {
	CMD    string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (client *Client) Command(command *Command) (err error) {
	session, err := client.SSH.NewSession()
	if err != nil {
		return
	}

	defer session.Close()

	term := os.Getenv("TERM")

	if term == "" {
		term = "xterm-256color"
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if command.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return err
		}

		go io.Copy(stdin, command.Stdin)
	}

	if command.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return err
		}

		go io.Copy(command.Stdout, stdout)
	}

	if command.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return err
		}

		go io.Copy(command.Stderr, stderr)
	}

	for _, env := range command.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if err = session.RequestPty(term, 40, 80, modes); err != nil {
		return
	}

	if err = session.Run(command.CMD); err != nil {
		return
	}

	return
}

func (client *Client) Shell() (err error) {
	session, err := client.SSH.NewSession()
	if err != nil {
		return
	}

	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	term := os.Getenv("TERM")

	if term == "" {
		term = "xterm"
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fd := int(os.Stdin.Fd())

	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return
	}

	defer terminal.Restore(fd, state)

	width, height, err := terminal.GetSize(fd)
	if err != nil {
		return
	}

	if err = session.RequestPty(term, height, width, modes); err != nil {
		return
	}

	if err = session.Shell(); err != nil {
		return
	}

	if err = session.Wait(); err != nil {
		return
	}

	return
}
