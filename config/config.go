package main

import (
	"io/ioutil"

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
	AppYamlPath    string `yaml:",omitempty"`
	PrivateKeyPath string `yaml:"private_key_path"`
	PublicKeyPath  string `yaml:"public_key_path"`
	VaultPath      string `yaml:"vault_path"`
	UserName       string `yaml:"user_name"`
}

func (c *Config) Set() error {
	if len(c.AppYamlPath) > 0 {

	} else {
		c.Version = VERSION
		c.Revision = REVISION
		c.UserName = ""
		c.PrivateKeyPath = "./cryptorious_privatekey"
		c.PublicKeyPath = "./cryptorious_publickey"
		c.VaultPath = "./cryptorious_vault.yaml"
		c.AppYamlPath = "./cryptorious.yaml"

		yamlBytes, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(c.AppYamlPath, yamlBytes, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Load() error {

	return nil
}

// Configuration() returns the configuration for application level logic
func GetConfiguration() (c Config, err error) {
	if err := c.Set(); err != nil {
		return c, err
	}
	return c, nil
}
