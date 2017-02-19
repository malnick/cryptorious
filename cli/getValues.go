package cli

import (
	"fmt"

	"github.com/malnick/cryptorious/action"
	gc "github.com/rthornton128/goncurses"
)

func vaultSetFromCurses() (*action.VaultSet, error) {
	vaultSet := &action.VaultSet{}

	username, err := getValuesFor("Username")
	if err != nil {
		return vaultSet, err
	}
	vaultSet.Username = username

	password, err := getValuesFor("Password")
	if err != nil {
		return vaultSet, err
	}
	vaultSet.Password = password

	note, err := getValuesFor("Secure Note")
	if err != nil {
		return vaultSet, err
	}
	vaultSet.SecureNote = note

	return vaultSet, nil
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
