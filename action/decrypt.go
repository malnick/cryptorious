package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
)

func Decrypt(key string, c config.Config) error {
	privData, err := ioutil.ReadFile(c.PrivateKeyPath)
	if err != nil {
		log.Errorf("%s was not found. Try `generate` first.", c.PrivateKeyPath)
		return err
	}
	log.Debug("Private key file: ", c.PrivateKeyPath)

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

	log.Debug("Retreiving encrypted values from vault...")
	username, encryptedPassword, encryptedNote, err := lookUpVault(key, c)
	if err != nil {
		return err
	}

	log.Debug("Decrypting password...")
	decryptedPassword, err := decryptValue(priv, encryptedPassword)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Decrypting notes...")
	decryptedNote, err := decryptValue(priv, encryptedNote)
	if err != nil {
		return err
	}

	printDecrypted(key, username, string(decryptedPassword), string(decryptedNote))

	return nil
}

func decryptValue(privkey *rsa.PrivateKey, encryptedValue string) ([]byte, error) {
	if encryptedValue == "" {
		log.Warn("Encrypted value empty, skipping.")
		return []byte(""), nil
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privkey, []byte(encryptedValue))
}

func lookUpVault(key string, c config.Config) (string, string, string, error) {
	var vault = vault.Vault{}
	vault.Path = c.VaultPath
	vault.Load()
	if _, ok := vault.Data[key]; !ok {
		return "", "", "", errors.New(fmt.Sprintf("%s not found in %s", key, vault.Path))
	}
	return vault.Data[key].Username, vault.Data[key].Password, vault.Data[key].SecureNote, nil
}
