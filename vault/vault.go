package vault

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Vault is the abstraction around the encrypted yaml vault type
type Vault struct {
	Data map[string]*Set `yaml:"data"`
	Path string
	Dir  string
}

// New returns a new vault, loaded from disk, at the given path.
func New(path string) (Vault, error) {
	v := Vault{
		Path: path,
	}
	err := v.Load()
	if err != nil {
		return v, err
	}
	return v, nil
}

// Load unmarshals the YAML from disk to a serialized object for CRUD operations.
func (v *Vault) Load() error {
	if _, err := os.Stat(v.Path); err != nil {
		log.Warnf("%s not found, can not load vault", v.Path)
		return nil
	}
	yamlBytes, err := ioutil.ReadFile(v.Path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlBytes, &v.Data)
	if err != nil {
		return err
	}
	return nil
}

// Write marshals the data to YAML and writes to disk.
// NOTE: assumes .Load() has been called.
func (v *Vault) Write() error {
	newYamlData, err := yaml.Marshal(&v.Data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(v.Path); err != nil {
		log.Warnf("%s does not exist, writing new vault file.", v.Path)
	}
	return ioutil.WriteFile(v.Path, newYamlData, 0644)
}

// Add creates a new vault entry
func (v *Vault) Add(key string, vs *Set) error {
	if v.Data == nil {
		v.Data = make(map[string]*Set)
	}

	if _, ok := v.Data[key]; ok {
		return fmt.Errorf("vault entry for %s found, try `cryptorious delete %s` first?", key, key)
	}

	log.Infof("adding new vault entry for %s", key)
	v.Data[key] = vs

	return v.Write()
}

// Delete removes an entry from the vault and writes the updated vault to disk.
// NOTE: Assums .Load() as been called.
func (v *Vault) Delete(key string) error {
	_, ok := v.Data[key]
	if !ok {
		return fmt.Errorf("Vault entry for %s not found, can not remove", key)
	}

	delete(v.Data, key)

	return v.Write()
}
