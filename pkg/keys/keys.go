package keys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// GetKeyPair will attempt to get the keypair from a file and will fail back
// to generating a new set and saving it to the file. Returns pub, priv, err
func ReadGenerateKeyPair(keyPairName string) (string, string, error) {
	// read keys from file
	pub, priv, err := ReadKeyPair(keyPairName)
	if err != nil {
		goto GENERATE_KEYS
	} else {
		return string(pub), string(priv), nil
	}

	// generate keys and save to file
GENERATE_KEYS:
	pub, priv, err = GenerateKeyPair()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate keys - %s", err)
	}

	if err = WriteKeyPair(keyPairName, pub, priv); err != nil {
		return "", "", fmt.Errorf("failed to write file - %s", err)
	}

	return pub, priv, nil
}

// GenerateKeyPair make a pair of public and private keys for SSH access.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
func GenerateKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	var private bytes.Buffer

	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	public := ssh.MarshalAuthorizedKey(pub)

	return string(public), private.String(), nil
}

func ReadKeyPair(keyPair string) (string, string, error) {
	_, err := os.Stat(keyPair)
	if err != nil {
		return "", "", err
	}

	priv, err := ioutil.ReadFile(keyPair)
	if err != nil {
		return "", "", err
	}

	pub, err := ioutil.ReadFile(keyPair + ".pub")
	if err != nil {
		return "", "", err
	}

	return string(pub), string(priv), nil
}

func WriteKeyPair(keyPairName, pub, priv string) error {
	directory := filepath.Dir(keyPairName)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if directory != "" {
			err = os.MkdirAll(directory, os.ModePerm)
			if err != nil {
				return err
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := ioutil.WriteFile(keyPairName, []byte(priv), 0600); err != nil {
		return fmt.Errorf("failed to write file - %s", err)
	}

	if err := ioutil.WriteFile(keyPairName+".pub", []byte(pub), 0644); err != nil {
		return fmt.Errorf("failed to write pub file - %s", err)
	}

	return nil
}
