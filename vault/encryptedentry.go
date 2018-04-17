package vault

import (
	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/aws/kms"
)

// EncryptedEntry contains the encrypted values of a vault set
type EncryptedEntry struct {
	Ciphertext string `yaml:"ciphertext"`
	IV         string `yaml:"iv"`
	Key        string `yaml:"key"`
}

// Set contains the values of a vault entry
type Set struct {
	Username   EncryptedEntry `yaml:"username"`
	Password   EncryptedEntry `yaml:"password"`
	SecureNote EncryptedEntry `yaml:"secure_note"`
}

// Decrypt returns the cleartext values of an encrypted vault set
func (vs *Set) Decrypt(kms kms.Impl) (CleartextEntry, error) {
	clr := CleartextEntry{}

	if len(vs.Username.Ciphertext) != 0 {
		log.Debug("decrypting username...")
		username, err := kms.Decrypt(
			[]byte(vs.Username.Ciphertext),
			[]byte(vs.Username.Key),
			[]byte(vs.Username.IV))
		if err != nil {
			return clr, err
		}
		log.Debugf("cleartext: %s", username)
		clr.Username = username
	}

	if len(vs.Password.Ciphertext) != 0 {
		log.Debug("decrypting password...")
		password, err := kms.Decrypt(
			[]byte(vs.Password.Ciphertext),
			[]byte(vs.Password.Key),
			[]byte(vs.Password.IV))
		if err != nil {
			return clr, err
		}
		log.Debugf("cleartext: %s", password)
		clr.Password = password
	}

	if len(vs.SecureNote.Ciphertext) != 0 {
		log.Debug("decrypting secure note...")
		secureNote, err := kms.Decrypt(
			[]byte(vs.SecureNote.Ciphertext),
			[]byte(vs.SecureNote.Key),
			[]byte(vs.SecureNote.IV))
		if err != nil {
			return clr, err
		}
		log.Debugf("cleartext: %s", secureNote)
		clr.SecureNote = secureNote
	}

	return clr, nil
}
