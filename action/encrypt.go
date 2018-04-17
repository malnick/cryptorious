package action

import (
	"fmt"

	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
	gc "github.com/rthornton128/goncurses"
)

// Encrypt accepts a key and cryptorious config and returns an error
// if found during the encryption process
func Encrypt(key string, c config.Config) error {
	thisVault, err := vault.New(c.VaultPath)
	if err != nil {
		return err
	}

	clr, err := cleartextFromCurses()
	if err != nil {
		return err
	}

	vs, err := clr.Encrypt(c.KMSClient, c.KMSKeyARN)
	if err != nil {
		return err
	}

	return thisVault.Add(key, vs)
}

func cleartextFromCurses() (*vault.CleartextEntry, error) {
	clr := &vault.CleartextEntry{}

	username, err := getValuesFor("Username")
	if err != nil {
		return clr, err
	}
	clr.Username = username

	password, err := getValuesFor("Password")
	if err != nil {
		return clr, err
	}
	clr.Password = password

	note, err := getValuesFor("Secure Note")
	if err != nil {
		return clr, err
	}
	clr.SecureNote = note

	return clr, nil
}

func getValuesFor(key string) (string, error) {
	stdscr, _ := gc.Init()
	defer gc.End()

	prompt := fmt.Sprintf("Enter %s: ", key)
	row, col := stdscr.MaxYX()
	row, col = (row/2)-1, (col-len(prompt))/2
	stdscr.MovePrint(row, col, prompt)

	/* GetString will only retieve the specified number of characters. Any
	attempts by the user to enter more characters will elicit an audiable
		beep */
	var value string
	value, err := stdscr.GetString(10000)
	if err != nil {
		return value, err
	}

	//	stdscr.Refresh()
	stdscr.GetChar()
	stdscr.Erase()

	return value, nil
}
