package config

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var (
	VERSION  string
	REVISION string
)

// Config{} is the complete application configuration
type Config struct {
	DebugMode      bool
	Version        string
	Revision       string
	PrivateKeyPath string
	PublicKeyPath  string
	VaultDir       string
	VaultPath      string
}

// set() configurations application level direcotories such as the .cryptorious $HOME dir, and .ssh if it does not exist.
func (c *Config) setDefaults() error {
	home := os.Getenv("HOME")
	c.DebugMode = false
	c.Version = VERSION
	c.Revision = REVISION
	c.PrivateKeyPath = fmt.Sprintf("%s/.ssh/cryptorious_privatekey", home)
	c.PublicKeyPath = fmt.Sprintf("%s/.ssh/cryptorious_publickey", home)
	c.VaultDir = fmt.Sprintf("%s/.cryptorious", home)
	c.VaultPath = fmt.Sprintf("%s/vault.yaml", c.VaultDir)
	if err := statDirectoryOrCreate(c.VaultDir); err != nil {
		return err
	}
	return nil
}

func statDirectoryOrCreate(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		log.Warnf("%s does not exist, creating...", dir)
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// Configuration() returns the configuration for application level logic
func GetConfiguration() (c Config, err error) {
	if err := c.setDefaults(); err != nil {
		return c, err
	}
	return c, nil
}
