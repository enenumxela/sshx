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

// RunCommand
func (client *Client) RunCommand(command *Command) (err error) {
	session, err := client.SSH.NewSession()
	if err != nil {
		return
	}

	defer session.Close()

	for _, env := range command.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
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

	command.CMD = "source ~/.profile && " + command.CMD

	if err = session.Run(command.CMD); err != nil {
		return
	}

	return
}

// GetShell
func (client *Client) GetShell() (err error) {
	session, err := client.SSH.NewSession()
	if err != nil {
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
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

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err = session.Shell(); err != nil {
		return
	}

	if err = session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}

		return
	}

	return
}
