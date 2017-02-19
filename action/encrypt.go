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

func Encrypt(key string, vs *vault.VaultSet, c config.Config) error {
	pubkey, err := createPublicKey(c.PublicKeyPath)
	if err != nil {
		return err
	}

	thisVault, err := vault.New(c.VaultPath)
	if err != nil {
		return err
	}

	if len(vs.Password) > 0 {
		if encoded, err := encryptValue(pubkey, vs.Password); err == nil {
			vs.Password = string(encoded)
		} else {
			return err
		}
	}

	if len(vs.SecureNote) > 0 {
		if encoded, err := encryptValue(pubkey, vs.SecureNote); err == nil {
			vs.SecureNote = string(encoded)
		} else {
			return err
		}
	}

	if len(vs.Username) > 0 {
		vs.Username = vs.Username
	}

	if err := thisVault.Update(key, vs); err != nil {
		return err
	}

	return nil
}

func encryptValue(pubkey interface{}, value string) ([]byte, error) {
	// Encode the passed in value
	log.Debugf("Encoding value: %s", value)
	encodedValue, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey.(*rsa.PublicKey), []byte(value))
	return encodedValue, err
}

func createPublicKey(path string) (*rsa.PublicKey, error) {
	pubData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	log.Debug("using public key file: ", path)
	log.Debug(string(pubData))

	pubkey, err := createPublicKeyBlockCipher(pubData)
	if err != nil {
		return nil, err
	}

	return pubkey.(*rsa.PublicKey), nil
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
