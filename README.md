# sshx

A [Go(Golang)](https://golang.org/) package to provide a simple abstraction around [ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) and [sftp](https://pkg.go.dev/github.com/pkg/sftp) packages.

## Resources

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
    * [Start Connection With Password](#start-connection-with-password)
    * [Start Connection with private key (With Passphrase)](#start-connection-with-private-key-with-passphrase)
    * [Start Connection With Private Key (Without Passphrase)](#start-connection-with-private-key-without-passphrase)
    * [Upload Local File to Remote](#upload-local-file-to-remote)
    * [Download Remote File to Local](#download-remote-file-to-local)
    * [Run Remote Commands](#run-remote-commands)
    * [Get an Interactive Shell](#get-an-interactive-shell)

## Features

* Easy to use and simple API.
* Supports connections with passwords.
* Supports connections with private keys (with passphrase).
* Supports connections with private keys (without passphrase).
* Supports upload files from local to remote.
* Supports download files from remote to local.
* Supports running remote commands.
* Supports getting an interactive shell.

## Installation

Install sshx as you normally would for any Go package:

```bash
go get -u github.com/enenumxela/sshx
```

## Usage

### Start Connection With Password

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

### Start Connection with Private Key (With Passphrase)

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

### Start Connection With Private Key (Without Passphrase)

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

### Run Remote Commands

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

### Get an Interactive Shell

```go
if err = client.Shell(); err != nil {
    log.Println(err)
}
```