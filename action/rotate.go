package action

import (
	"github.com/malnick/cryptorious/config"
)

// 1. Backup old keys
// 2. Generate new keys
// 3. Load vault into memory and decrypt all passwords and notes with old keys (now at backup path)
// 4. Encrypt vault passwords and notes with new keys
// 5. Write vault to disk

// RotateVault is currently not implemented
func RotateVault(c config.Config) error {
	// not implemented

	return nil
}
