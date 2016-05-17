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
func GenerateKeys(c config.Config) {
	keyPath := c.KeyDirectory
	username := c.Username
	privPath := fmt.Sprintf("%s/privatekey_%s.pem", keyPath, username)
	pubPath := fmt.Sprintf("%s/publickey_%s.pem", keyPath, username)
	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)
	checkError(err)
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
	checkError(err)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: ansipub,
	})
	ioutil.WriteFile(pubPath, pubBytes, 0644)
	log.Info("Public Key: ", pubPath)
	fmt.Println(string(pubBytes))
}
