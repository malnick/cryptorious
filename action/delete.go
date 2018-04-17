package action

import (
	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/vault"
)

// DeleteVaultEntry remove a key and its values from the vault
func DeleteVaultEntry(key string, vaultPath string) error {
	vault, err := vault.New(vaultPath)
	if err != nil {
		return err
	}

	log.Warnf("Removing '%s' entry from vault", key)
	return vault.Delete(key)
}
