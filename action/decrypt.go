package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
	"github.com/olekukonko/tablewriter"
)

func Decrypt(key string, c config.Config) error {
	privData, err := ioutil.ReadFile(c.PrivateKeyPath)
	if err != nil {
		log.Errorf("%s was not found. Try `generate` first.", c.PrivateKeyPath)
		return err
	}
	log.Debug("Private key file: ", c.PrivateKeyPath)
	log.Debug(string(privData))
	if err != nil {
		return err
	}
	// Extract the PEM-encoded data block
	block, _ := pem.Decode(privData)
	if block == nil {
		log.Error("bad key data: %s", "not PEM-encoded")
		return err
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Error("unknown key type %q, want %q", got, want)
		return err
	}
	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error("bad private key: %s", err)
		return err
	}

	username, encryptedPassword, encryptedNote, err := lookUpVault(key, c)
	if err != nil {
		return err
	}

	decryptedPassword, err := decryptValue(priv, encryptedPassword)
	if err != nil {
		return err
	}

	decryptedNote, err := decryptValue(priv, encryptedNote)
	if err != nil {
		return err
	}

	prettyPrintMe := [][]string{
		[]string{key, username, string(decryptedPassword), string(decryptedNote)},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Username", "Password", "Secure Note"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prettyPrintMe) // Add Bulk Data
	table.Render()

	return nil
}

func decryptValue(privkey *rsa.PrivateKey, encryptedValue string) ([]byte, error) {
	return rsa.DecryptOAEP(sha1.New(), rand.Reader, privkey, []byte(encryptedValue), []byte(">"))
}

func lookUpVault(key string, c config.Config) (string, string, string, error) {
	var vault = Vault{}
	vault.Path = c.VaultPath
	vault.load()
	if _, ok := vault.Data[key]; !ok {
		return "", "", "", errors.New(fmt.Sprintf("%s not found in %s", key, vault.Path))
	}
	return vault.Data[key].Username, vault.Data[key].Password, vault.Data[key].SecureNote, nil
}
