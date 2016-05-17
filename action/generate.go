package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
)

// GenerateKeys() creates public private keys for a $USER
func GenerateKeys(c config.Config) error {
	privPath := c.PrivateKeyPath
	pubPath := c.PublicKeyPath
	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	// Write Private Key
	privBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	})
	ioutil.WriteFile(privPath, privBytes, 0600)
	log.Info("Private Key: ", privPath)
	fmt.Println(string(privBytes))
	// Write Public Key
	ansipub, err := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
	if err != nil {
		return err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: ansipub,
	})
	ioutil.WriteFile(pubPath, pubBytes, 0644)
	log.Info("Public Key: ", pubPath)
	fmt.Println(string(pubBytes))
	return nil
}
