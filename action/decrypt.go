package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
)

func Decrypt(key string, c config.Config) error {
	privData, err := ioutil.ReadFile(c.PrivateKeyPath)
	if err != nil {
		log.Errorf("%s was not found. Try `generate` first.", c.PrivateKeyPath)
		return err
	}
	log.Debug("Private key file: ", c.PrivateKeyPath)
	log.Debug(string(privData))
	if err != nil {
		return err
	}
	// Extract the PEM-encoded data block
	block, _ := pem.Decode(privData)
	if block == nil {
		log.Error("bad key data: %s", "not PEM-encoded")
		return err
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Error("unknown key type %q, want %q", got, want)
		return err
	}
	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error("bad private key: %s", err)
		return err
	}

	encryptedValue, err := lookUpValueFromVault(key, c)
	if err != nil {
		return err
	}

	log.Debugf("%s found in %s", key, c.VaultPath)

	decryptedValue, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, priv, []byte(encryptedValue), []byte(">"))
	if err != nil {
		return err
	}

	fmt.Printf("Decrypted value for %s => %s\n", key, decryptedValue)

	return nil
}

func lookUpValueFromVault(key string, c config.Config) (string, error) {
	var vault = Vault{}
	vault.Path = c.VaultPath
	vault.load()
	if _, ok := vault.Data[key]; !ok {
		return "", errors.New(fmt.Sprintf("%s not found in %s", key, vault.Path))
	}
	return vault.Data[key].Password, nil
}
