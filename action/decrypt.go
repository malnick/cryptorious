package action

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
	"github.com/olekukonko/tablewriter"
)

func Decrypt(key string, c config.Config) error {
	log.Debug("Reading key file ", c.KeyPath)
	keydata, err := ioutil.ReadFile(c.KeyPath)
	if err != nil {
		log.Errorf("%s was not found. Try `generate` first.", c.KeyPath)
		return err
	}

	username, encryptedPassword, encryptedNote, err := lookUpVault(key, c)
	if err != nil {
		return err
	}

	decryptedPassword, err := decryptValue(keydata, []byte(encryptedPassword))
	if err != nil {
		return err
	}

	decryptedNote, err := decryptValue(keydata, []byte(encryptedNote))
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

func decryptValue(key, ciphertext []byte) ([]byte, error) {
	var block cipher.Block

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	if len(ciphertext) < aes.BlockSize {
		err := errors.New("ciphertext too short")
		return []byte{}, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	plaintext := ciphertext
	return plaintext, nil
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
