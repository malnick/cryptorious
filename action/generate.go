package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// GenerateKeys creates public private keys for a $USER
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

	if checkFileExists(privPath) {
		log.Warnf("%s exists, please manually remove to proceed.", privPath)
		return errors.New("Will not overwrite existing private key path.")
	}
	if err := ioutil.WriteFile(privPath, privBytes, 0600); err != nil {
		return err
	}

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

	if checkFileExists(pubPath) {
		log.Warnf("%s exists, please manually remove to proceed.", pubPath)
		return errors.New("Will not overwrite existing public key path.")
	}
	if err := ioutil.WriteFile(pubPath, pubBytes, 0644); err != nil {
		return err
	}
	log.Info("Public Key: ", pubPath)
	fmt.Println(string(pubBytes))
	return nil
}

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

func NewPassword(length int) error {
	p, err := randomPassword(length, StdChars)
	fmt.Println(p)
	return err

}

func randomPassword(length int, chars []byte) (string, error) {
	// https://raw.githubusercontent.com/cmiceli/password-generator-go/master/gen.go

	newPassword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			return string(newPassword), err
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPassword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPassword), nil
			}
		}
	}

}
