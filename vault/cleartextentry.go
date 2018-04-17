package vault

import (
	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/aws/kms"
)

// CleartextEntry holds the cleartext values of a vault set
type CleartextEntry struct {
	Username   string
	Password   string
	SecureNote string
}

// Encrypt returns an encrypted vault set for a cleartext value
func (c *CleartextEntry) Encrypt(kms kms.Impl, keyARN string) (*Set, error) {
	vaultSet := &Set{
		Username:   EncryptedEntry{},
		Password:   EncryptedEntry{},
		SecureNote: EncryptedEntry{},
	}

	log.Debug("encrypting username...")
	if len(c.Username) != 0 {
		ciphertext, key, iv, err := kms.Encrypt([]byte(c.Username), keyARN)
		if err != nil {
			return vaultSet, err
		}

		vaultSet.Username.Ciphertext = string(ciphertext)
		log.Debugf("ciphertext: %s", ciphertext)
		vaultSet.Username.IV = string(iv)
		log.Debugf("iv: %s", iv)
		vaultSet.Username.Key = string(key)
		log.Debugf("key: %s", key)
	}

	log.Debug("encrypting password...")
	if len(c.Password) != 0 {
		ciphertext, key, iv, err := kms.Encrypt([]byte(c.Password), keyARN)
		if err != nil {
			return vaultSet, err
		}

		vaultSet.Password.Ciphertext = string(ciphertext)
		log.Debugf("ciphertext: %s", ciphertext)
		vaultSet.Password.IV = string(iv)
		log.Debugf("iv: %s", iv)
		vaultSet.Password.Key = string(key)
		log.Debugf("key: %s", key)
	}

	log.Debug("encrypting secure note...")
	if len(c.SecureNote) != 0 {
		ciphertext, key, iv, err := kms.Encrypt([]byte(c.SecureNote), keyARN)
		if err != nil {
			return vaultSet, err
		}

		vaultSet.SecureNote.Ciphertext = string(ciphertext)
		log.Debugf("ciphertext: %s", ciphertext)
		vaultSet.SecureNote.IV = string(iv)
		log.Debugf("iv: %s", iv)
		vaultSet.SecureNote.Key = string(key)
		log.Debugf("key: %s", key)

	}

	return vaultSet, nil
}
