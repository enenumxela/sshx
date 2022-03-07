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

	for _, env := range command.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")

	if term == "" {
		term = "xterm-256color"
	}

	if err = session.RequestPty(term, 40, 80, modes); err != nil {
		return
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

	return session.Run(command.CMD)
}

func (client *Client) Shell() (err error) {
	session, err := client.SSH.NewSession()
	if err != nil {
		return
	}

	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")

	if term == "" {
		term = "xterm-256color"
	}

	fd := int(os.Stdin.Fd())

	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return
	}

	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		return
	}

	if err = session.RequestPty(term, h, w, modes); err != nil {
		return
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err = session.Shell(); err != nil {
		return
	}

	return session.Wait()
}
