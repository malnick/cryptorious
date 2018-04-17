package config

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/aws"
	"github.com/malnick/cryptorious/aws/kms"
)

var (
	// VERSION is the git tag
	VERSION string
	// REVISION is the short sha of the commit
	REVISION string
)

// Config is the complete application configuration
type Config struct {
	DebugMode             bool
	Version               string
	Revision              string
	KMSKeyARN             string
	KMSClient             kms.Impl
	VaultDir              string
	VaultPath             string
	Clipboard             bool
	Goto                  bool
	DecryptSessionTimeout int
	PrintAll              bool
}

// setDefaults configurations application level direcotories such as the .cryptorious $HOME dir, and .ssh if it does not exist.
func (c *Config) setDefaults() error {
	home := os.Getenv("HOME")
	c.DebugMode = false
	c.Version = VERSION
	c.Revision = REVISION
	c.VaultDir = fmt.Sprintf("%s/.cryptorious", home)
	c.VaultPath = fmt.Sprintf("%s/vault.yaml", c.VaultDir)

	a, _ := aws.New()
	c.KMSClient = kms.New(a)

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

// GetConfiguration returns the configuration for application level logic
func GetConfiguration() (c Config, err error) {
	if err := c.setDefaults(); err != nil {
		return c, err
	}
	return c, nil
}
