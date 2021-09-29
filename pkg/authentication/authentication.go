package authentication

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// Authentication represents ssh auth methods.
type Authentication []ssh.AuthMethod

// Password returns password auth method.
func Password(password string) Authentication {
	return Authentication{
		ssh.Password(password),
	}
}

// Key returns auth method from private key with or without passphrase.
func Key(prvFile string, passphrase string) (Authentication, error) {

	signer, err := GetSigner(prvFile, passphrase)

	if err != nil {
		return nil, err
	}

	return Authentication{
		ssh.PublicKeys(signer),
	}, nil
}

// GetSigner returns ssh signer from private key file.
func GetSigner(prvFile string, passphrase string) (ssh.Signer, error) {

	var (
		err    error
		signer ssh.Signer
	)

	privateKey, err := ioutil.ReadFile(prvFile)

	if err != nil {

		return nil, err

	} else if passphrase != "" {

		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))

	} else {

		signer, err = ssh.ParsePrivateKey(privateKey)
	}

	return signer, err
}
