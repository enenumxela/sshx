package sshx

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/enenumxela/sshx/pkg/authentication"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// Client is the structure of sshx client.
type Client struct {
	SSH  *ssh.Client
	SFTP *sftp.Client
}

// Configuration is the structure of configuration for sshx client.
type Configuration struct {
	Auth            authentication.Authentication
	User            string
	Addr            string
	Port            uint
	Timeout         int
	HostKeyCallback ssh.HostKeyCallback
}

type Command struct {
	CMD    string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New create a new sshx NewClient
func New(configuration *Configuration) (client *Client, err error) {
	client = &Client{
		SSH:  &ssh.Client{},
		SFTP: &sftp.Client{},
	}

	if configuration.Timeout == 0 {
		configuration.Timeout = 20
	}

	client.SSH, err = ssh.Dial("tcp", net.JoinHostPort(configuration.Addr, fmt.Sprint(configuration.Port)), &ssh.ClientConfig{
		User:            configuration.User,
		Auth:            configuration.Auth,
		Timeout:         time.Duration(configuration.Timeout) * time.Second,
		HostKeyCallback: configuration.HostKeyCallback,
	})
	if err != nil {
		return
	}

	client.SFTP, err = sftp.NewClient(client.SSH)
	if err != nil {
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

func (client *Client) RunCommand(command *Command) (err error) {
	session, err := client.SSH.NewSession()
	if err != nil {
		return
	}

	defer session.Close()

	if err = client.PrepareCommand(session, command); err != nil {
		log.Fatal(err)
	}

	if err = session.Run(command.CMD); err != nil {
		return
	}

	return
}

func (client *Client) PrepareCommand(session *ssh.Session, cmd *Command) (err error) {
	cmd.CMD = "source ~/.profile && " + cmd.CMD

	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return err
		}

		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return err
		}

		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return err
		}

		go io.Copy(cmd.Stderr, stderr)
	}

	return
}

func (client *Client) Close() (err error) {
	if client == nil {
		return
	}

	if client.SFTP != nil {
		if err = client.SFTP.Close(); err != nil {
			return
		}
	}

	if client.SSH != nil {
		if err = client.SSH.Close(); err != nil {
			return
		}
	}

	return
}
