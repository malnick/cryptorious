package action

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/vault"
)

func RenameVaultEntry(oldkey, newkey, vaultPath string) error {
	v, err := vault.New(vaultPath)
	if err != nil {
		return err
	}

	if err := v.Load(); err != nil {
		return err
	}

	if _, ok := v.Data[oldkey]; !ok {
		return errors.New(fmt.Sprintf("%s is not a valid key in the vault", oldkey))
	}

	v.Data[newkey] = v.Data[oldkey]
	delete(v.Data, oldkey)

	log.Infof("%s -> %s, new vault updates written to disk.", oldkey, newkey)
	v.Write()

	return nil
}
