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
	"github.com/atotto/clipboard"
	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
)

func Decrypt(key string, c config.Config) error {
	priv, err := createPrivateKey(c.PrivateKeyPath)
	if err != nil {
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

	if c.Clipboard {
		log.Info("Copying decrypted password to clipboard!")
		return clipboard.WriteAll(string(decryptedPassword))
	}

	printDecrypted(key, username, string(decryptedPassword), string(decryptedNote))

	return nil
}

func createPrivateKey(path string) (*rsa.PrivateKey, error) {
	privData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("%s was not found. Try `generate` first.", path)
		return nil, err
	}
	log.Debug("Private key file: ", path)

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(privData)
	if block == nil {
		log.Error("bad key data: %s", "not PEM-encoded")
		return nil, err
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Error("unknown key type %q, want %q", got, want)
		return nil, err
	}
	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error("bad private key: %s", err)
		return nil, err
	}
	return priv, nil
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
