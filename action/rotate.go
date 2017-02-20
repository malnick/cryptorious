package action

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
)

// 1. Backup old keys
// 2. Generate new keys
// 3. Load vault into memory and decrypt all passwords and notes with old keys (now at backup path)
// 4. Encrypt vault passwords and notes with new keys
// 5. Write vault to disk

func RotateVault(c config.Config) error {
	// Backup old keys
	if err := backupKeys(c.PrivateKeyPath, c.PublicKeyPath); err != nil {
		return err
	}

	// Generate new keys
	if err := GenerateKeys(c); err != nil {
		return err
	}

	// Rotate vault data
	if err := rotateVaultData(c); err != nil {
		return err
	}

	return nil
}

func backupKeys(private, public string) error {
	newPriv := private + ".bak"
	newPub := public + ".bak"

	if err := os.Rename(private, newPriv); err != nil {
		return err
	}

	if err := os.Rename(public, newPub); err != nil {
		return err
	}

	return nil
}

func rotateVaultData(c config.Config) error {
	// Use the old key which has already been backed up for decryption
	privBk := c.PrivateKeyPath + ".bak"
	decryptKey, err := createPrivateKey(privBk)
	if err != nil {
		return err
	}

	// Use the new key, at the confing path, for encrypting.
	encryptKey, err := createPublicKey(c.PublicKeyPath)
	if err != nil {
		return err
	}

	v, err := vault.New(c.VaultPath)
	if err != nil {
		return err
	}

	// Backup old vault file
	if err := os.Rename(c.VaultPath, c.VaultPath+".bk"); err != nil {
		return err
	}

	// Decrypt the entire vault in place, re-encrypting along the way
	for key, set := range v.Data {
		log.Infof("Rotating entry for %s", key)
		decryptedPassword, err := decryptValue(decryptKey, set.Password)
		if err != nil {
			return err
		}

		encryptedPassword, err := encryptValue(encryptKey, string(decryptedPassword))
		if err != nil {
			return err
		}
		v.Data[key].Password = string(encryptedPassword)

		decryptedNote, err := decryptValue(decryptKey, set.SecureNote)
		if err != nil {
			return err
		}

		encryptedNote, err := encryptValue(encryptKey, string(decryptedNote))
		if err != nil {
			return err
		}
		v.Data[key].SecureNote = string(encryptedNote)
	}

	if err := v.Write(); err != nil {
		return err
	}

	return nil
}
