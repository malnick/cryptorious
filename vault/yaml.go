package vault

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type VaultSet struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	SecureNote string `yaml:"secure_note"`
}

type Vault struct {
	Data map[string]*VaultSet `yaml:"data"`
	Path string
	Dir  string
}

// New() returns a new vault, loaded from disk, at the given path.
func New(path string) (Vault, error) {
	vault := Vault{
		Path: path,
	}
	err := vault.Load()
	if err != nil {
		return vault, err
	}
	return vault, nil
}

// Load() unmarshals the YAML from disk to a serialized object for CRUD operations.
func (vault *Vault) Load() error {
	if _, err := os.Stat(vault.Path); err != nil {
		log.Warnf("%s not found, will create new Vault file.", vault.Path)
		return nil
	}
	yamlBytes, err := ioutil.ReadFile(vault.Path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlBytes, &vault.Data)
	if err != nil {
		return err
	}
	return nil
}

// Write() marshals the data to YAML and writes to disk.
// NOTE: assumes .Load() has been called.
func (vault *Vault) Write() error {
	newYamlData, err := yaml.Marshal(&vault.Data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(vault.Path); err != nil {
		log.Warnf("%s does not exist, writing new vault file.", vault.Path)
	}
	return ioutil.WriteFile(vault.Path, newYamlData, 0644)
}

// Delete() removes an entry from the vault and writes the updated vault to disk.
// NOTE: Assums .Load() as been called.
func (vault *Vault) Delete(key string) error {
	_, ok := vault.Data[key]
	if !ok {
		return errors.New(fmt.Sprintf("Vault entry for %s not found, can not remove", key))
	}

	delete(vault.Data, key)

	return vault.Write()
}
