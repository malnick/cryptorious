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
	"os"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
)

type Vault struct {
	Data map[string]string
	Path string
	Dir  string
}

func (vault *Vault) load() error {
	if _, err := os.Stat(vault.Path); err != nil {
		log.Warnf("%s not found, will create new Vault file.", vault.Path)
		return nil
	}
	yamlBytes, err := ioutil.ReadFile(vault.Path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlBytes, &vault.Data)
	if err != nil {
		return err
	}
	return nil
}

func (vault *Vault) writeValueToVault(key string, encodedValue []byte) error {
	// Assumes .load() was called before executing.
	newYamlData, err := yaml.Marshal(&vault.Data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(vault.Path); err != nil {
		log.Warnf("%s does not exist, writing new vault file.", vault.Path)
	}
	if err := ioutil.WriteFile(vault.Path, newYamlData, 0644); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"Key":          key,
		"[]Byte Value": encodedValue,
	}).Infof("Successfully wrote to %s", vault.Path)
	return nil
}

func Encrypt(key string, value string, c config.Config) error {
	pubData, err := ioutil.ReadFile(c.PublicKeyPath)
	if err != nil {
		return err
	}
	log.Debug("Using public key file: ", c.PublicKeyPath)
	log.Debug(string(pubData))

	pubkey, err := createPublicKeyBlockCipher(pubData)
	if err != nil {
		return err
	}

	// Encode the passed in value
	encodedValue, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pubkey.(*rsa.PublicKey), []byte(value), []byte(string(">")))
	if err != nil {
		return err
	}

	// Amend the Vault with the new data
	var vault = Vault{
		Data: make(map[string]string),
		Path: c.VaultPath,
	}

	vault.Data[key] = string(encodedValue)
	if err := vault.load(); err != nil {
		return err
	}

	if err := vault.writeValueToVault(key, encodedValue); err != nil {
		return err
	}

	return nil
}

func createPublicKeyBlockCipher(pubData []byte) (interface{}, error) {
	// Create block cipher from RSA key
	block, _ := pem.Decode(pubData)
	// Ensure key is PEM encoded
	if block == nil {
		return nil, errors.New(fmt.Sprintf("Bad key data: %s, not PEM encoded", string(pubData)))
	}
	// Ensure this is actually a RSA pub key
	if got, want := block.Type, "RSA PUBLIC KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("Unknown key type %q, want %q", got, want))
	}
	// Lastly, create the public key using the new block
	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}
