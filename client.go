package sshx

import (
	"fmt"
	"net"
	"time"

	"github.com/enenumxela/sshx/authentication"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client is the structure of sshx client.
type Client struct {
	SSH  *ssh.Client
	SFTP *sftp.Client
}

// Configuration is the structure of configuration for sshx client.
type Configuration struct {
	Addr            string
	Port            uint
	User            string
	Timeout         int
	Auth            authentication.Authentication
	HostKeyCallback ssh.HostKeyCallback
}

// New create a new sshx NewClient
func New(configuration *Configuration) (client *Client, err error) {
	client = &Client{
		SSH:  &ssh.Client{},
		SFTP: &sftp.Client{},
	}

	if client.SSH, err = ssh.Dial("tcp", net.JoinHostPort(configuration.Addr, fmt.Sprint(configuration.Port)), &ssh.ClientConfig{
		User:            configuration.User,
		Auth:            configuration.Auth,
		Timeout:         time.Duration(configuration.Timeout) * time.Second,
		HostKeyCallback: configuration.HostKeyCallback,
	}); err != nil {
		return
	}

	if client.SFTP, err = sftp.NewClient(client.SSH); err != nil {
		return
	}

	return
}

// Close closes the underlying client network connection.
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
