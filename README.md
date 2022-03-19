# sshx

A [Go(Golang)](https://golang.org/) package to provide a simple abstraction around [ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) and [sftp](https://pkg.go.dev/github.com/pkg/sftp) packages.

## Resources

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
    * [Connection with password](#connection-with-password)
    * [Connection with private key (with passphrase)](#connection-with-private-key-with-passphrase)
    * [Connection with private Key (without passphrase)](#connection-with-private-key-without-passphrase)
    * [Run remote commands](#run-remote-commands)
    * [Get an interactive shell](#get-an-interactive-shell)
    * [Upload Local File to Remote](#upload-local-file-to-remote)
    * [Download Remote File to Local](#download-remote-file-to-local)

## Features

* Supports connections with passwords.
* Supports connections with private keys (with passphrase).
* Supports connections with private keys (without passphrase).
* Supports running remote commands.
* Supports getting an interactive shell.
* Supports upload files from local to remote.
* Supports download files from remote to local.

## Installation

Install [sshx](https://github.com/enenumxela/sshx) as you normally would for any Go package:

```bash
go get -v -u github.com/enenumxela/sshx
```

## Usage

### Connection with password

```go
authentication, err := authentication.Password("Password")
if err != nil {
    log.Fatal(err)
}

client, err := sshx.New(&sshx.Configuration{
    Port:            22,
    Auth:            authentication,
    Addr:            "xxx.xxx.xxx.xxx",
    User:            "some-user",
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})
if err != nil {
    log.Println(err)
}

defer client.Close()
```

### Connection with private key (with passphrase)

```go
authentication, err := authentication.Key(privateKey, "Passphrase")
if err != nil {
    log.Fatal(err)
}

client, err := sshx.New(&sshx.Configuration{
    Port:            22,
    Auth:            authentication,
    Addr:            "xxx.xxx.xxx.xxx",
    User:            "some-user",
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})
if err != nil {
    log.Println(err)
}

defer client.Close()
```

### Connection with private key (without passphrase)

```go
authentication, err := authentication.Key(privateKey, "")
if err != nil {
    log.Fatal(err)
}

client, err := sshx.New(&sshx.Configuration{
    Port:            22,
    Auth:            authentication,
    Addr:            "xxx.xxx.xxx.xxx",
    User:            "some-user",
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})
if err != nil {
    log.Println(err)
}

defer client.Close()
```

### Run remote commands

```go
if err = client.Command(&sshx.Command{
    CMD:    "echo ${LC_TEST}",
    Env:    []string{"LC_TEST=working"},
    Stdin:  os.Stdin,
    Stdout: os.Stdout,
    Stderr: os.Stderr,
}); err != nil {
    log.Fatal(err)
}
```

### Get an interactive shell

```go
if err = client.Shell(); err != nil {
    log.Println(err)
}
```

### Upload Local File to Remote

```go
if err := client.Upload("/path/to/local/file", "/path/to/remote/file"); err != nil {
    log.Println(err)
}
```

### Download Remote File to Local

```go
if err := client.Download("/path/to/remote/file", "/path/to/local/file"); err != nil {
    log.Println(err)
}
```