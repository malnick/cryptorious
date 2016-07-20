package action

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// GenerateKeys creates public private keys for a $USER
func GenerateKeys(c config.Config) error {
	key := make([]byte, 8)
	keypath := c.KeyPath
	if _, err := rand.Read(key); err != nil {
		return err
	}

	if checkFileExists(keypath) {
		log.Warnf("%s exists, please manually remove to proceed.", keypath)
		return errors.New("Will not overwrite existing key.")
	}

	if err := ioutil.WriteFile(keypath, key, 0644); err != nil {
		return err
	}

	log.Info("Successfully wrote new AES key %s", keypath)
	fmt.Println(string(key))
	return nil
}
