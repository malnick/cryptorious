package action

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/atotto/clipboard"
	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
	"github.com/skratchdot/open-golang/open"
)

// Decrypt accepts a key and cryptorioius config and returns an error if found
// during the decryption process
func Decrypt(key string, c config.Config) error {
	log.Debug("Retreiving encrypted values from vault...")
	vs, err := lookUpVaultSet(key, c)
	if err != nil {
		return err
	}

	clr, err := vs.Decrypt(c.KMSClient)
	if err != nil {
		return err
	}

	if c.Clipboard {
		log.Info("Copying decrypted password to clipboard!")
		if err := clipboard.WriteAll(clr.Password); err != nil {
			return err
		}
	}

	if c.Goto {
		log.Infof("Opening default browser and logging into https://%s", key)
		if err := open.Run(fmt.Sprintf("https://%s", key)); err != nil {
			return err
		}
	}

	printDecrypted(key, clr.Username, clr.Password, clr.SecureNote, c.DecryptSessionTimeout)

	return nil
}

func lookUpVaultSet(key string, c config.Config) (*vault.Set, error) {
	var vault = vault.Vault{}
	vault.Path = c.VaultPath
	vault.Load()
	if _, ok := vault.Data[key]; !ok {
		return nil, fmt.Errorf("%s not found in %s", key, vault.Path)
	}
	return vault.Data[key], nil
}
