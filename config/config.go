package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	VERSION  string
	REVISION string
)

// Config{} is the complete application configuration
type Config struct {
	Version        string `yaml:",omitempty"`
	Revision       string `yaml:",omitempty"`
	AppYamlName    string `yaml:",omitempty"`
	VaultPath      string `yaml:",omitempty"`
	PrivateKeyName string `yaml:"private_key_path"`
	PublicKeyName  string `yaml:"public_key_path"`
	VaultName      string `yaml:"vault_path"`
	UserName       string `yaml:"user_name"`
}

func (c *Config) set() error {
	home := os.Getenv("HOME")
	if len(os.Getenv("HOME")) > 0 {
		c.VaultPath = fmt.Sprintf("%s/.cryptorious", home)
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", c.VaultPath, c.AppYamlName))
		if err != nil {
			if err := yaml.Unmarshal(fileBytes, &c); err != nil {
				return err
			}
		} else {
			c.Version = VERSION
			c.Revision = REVISION
			c.UserName = ""
			c.PrivateKeyName = "cryptorious_privatekey"
			c.PublicKeyName = "cryptorious_publickey"
			c.VaultName = "cryptorious_vault.yaml"
			c.AppYamlName = "cryptorious.yaml"

			log.Warn("Config file not found, writing one with all defaults ", c.AppYamlName)

			yamlBytes, err := yaml.Marshal(&c)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(c.AppYamlName, yamlBytes, 0644); err != nil {
				return err
			}
		}
	}
	return errors.New("Could not find $HOME directory, please make sure this environment variable is set before proceeding.")
}

// Configuration() returns the configuration for application level logic
func GetConfiguration() (c Config, err error) {
	if err := c.set(); err != nil {
		return c, err
	}
	return c, nil
}
