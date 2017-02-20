package action

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
	gc "github.com/rthornton128/goncurses"
)

func PrintAll(c config.Config) error {
	log.Info("Opening session for all entries in vault")
	v, err := getDecryptedVault(c)
	if err != nil {
		return err
	}

	menuItems := []string{}
	for entryName, _ := range v.Data {
		menuItems = append(menuItems, entryName)
	}

	printWithMenu(menuItems, v)
	return nil
}

func getDecryptedVault(c config.Config) (vault.Vault, error) {
	v, err := vault.New(c.VaultPath)
	if err != nil {
		return v, err
	}

	key, err := createPrivateKey(c.PrivateKeyPath)
	if err != nil {
		return v, err
	}
	if err := v.Load(); err != nil {
		return v, err
	}
	for name, _ := range v.Data {
		decryptedPassword, err := decryptValue(key, v.Data[name].Password)
		if err != nil {
			return v, err
		}
		v.Data[name].Password = string(decryptedPassword)

		decryptedNote, err := decryptValue(key, v.Data[name].SecureNote)
		if err != nil {
			return v, err
		}
		v.Data[name].SecureNote = string(decryptedNote)
	}

	return v, nil
}

func printWithMenu(menuItems []string, v vault.Vault) error {
	defer gc.End()
	stdscr := getDefaultScreen()

	gc.StartColor()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(3, gc.C_MAGENTA, gc.C_BLACK)

	items := make([]*gc.MenuItem, len(menuItems))
	for i, val := range menuItems {
		items[i], _ = gc.NewItem(val, "")
		defer items[i].Free()

		if i == 2 || i == 4 {
			items[i].Selectable(false)
		}
	}

	// create the menu
	menu, _ := gc.NewMenu(items)
	defer menu.Free()

	y, _ := stdscr.MaxYX()
	stdscr.MovePrint(y-3, 0, "Use up/down arrows to move; 'q' to exit")
	stdscr.Refresh()

	menu.SetForeground(gc.ColorPair(1) | gc.A_REVERSE)
	menu.SetBackground(gc.ColorPair(2) | gc.A_BOLD)
	menu.Grey(gc.ColorPair(3) | gc.A_BOLD)

	menu.Post()
	defer menu.UnPost()

	for {
		gc.Update()
		ch := stdscr.GetChar()
		switch ch {
		case ' ':
			menu.Driver(gc.REQ_TOGGLE)
		case 'q':
			return nil
		case gc.KEY_RETURN:
			stdscr.Move(20, 0)
			stdscr.ClearToEOL()

			entry := menu.Current(nil).Name()
			password := v.Data[entry].Password
			username := v.Data[entry].Username
			note := v.Data[entry].SecureNote
			printDecrypted(entry, username, password, note, 5)

			stdscr.Printf("Item selected is: %s", menu.Current(nil).Name())
			menu.PositionCursor()
		default:
			menu.Driver(gc.DriverActions[ch])
		}
	}
}

func printDecrypted(key, username, password, note string, sessionTimeout int) {
	defer gc.End()
	stdscr := getDefaultScreen()

	if !gc.HasColors() {
		log.Fatal("Cryptorious requires a colour capable terminal")
	}

	// Must be called after Init but before using any colour related functions
	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(3, gc.C_BLUE, gc.C_BLACK)

	stdscr.ColorOn(1)

	stdscr.AttrOn(gc.A_BOLD)
	stdscr.Border(gc.ACS_VLINE, gc.ACS_VLINE, gc.ACS_HLINE, gc.ACS_HLINE,
		gc.ACS_ULCORNER, gc.ACS_URCORNER, gc.ACS_LLCORNER, gc.ACS_LRCORNER)
	stdscr.AttrOff(gc.A_BOLD)

	row, col := stdscr.MaxYX()

	// Print the title of the table
	stdscr.ColorOn(2)

	keyMsg := "Decrypted Vault Entry for: "
	stdscr.HLine((row/2)-3, (col/2)-len(keyMsg), gc.ACS_HLINE, col/2)
	stdscr.HLine((row/2)-1, (col/2)-len(keyMsg), gc.ACS_HLINE, col/2)
	stdscr.HLine((row/2)+1, (col/2)-len(keyMsg), gc.ACS_HLINE, col/2)

	sep := make([]string, len(keyMsg))
	for i, _ := range sep {
		sep[i] = "-"
	}

	stdscr.MovePrint((row/2)-4, (col/2)-len(keyMsg), keyMsg)

	stdscr.ColorOn(3)
	stdscr.MovePrint((row/2)-4, (col / 2), fmt.Sprintf("'%s'", key))

	// Print the table keys
	stdscr.ColorOn(2)
	stdscr.AttrOn(gc.A_DIM)
	stdscr.MovePrint((row/2)-2, (col/2)-len(keyMsg), "Username: ")
	stdscr.MovePrint(row/2, (col/2)-len(keyMsg), "Password: ")
	stdscr.MovePrint((row/2)+2, (col/2)-len(keyMsg), "Secure Note: ")
	stdscr.AttrOff(gc.A_DIM)

	// Print the table values
	stdscr.ColorOn(1)
	stdscr.MovePrint((row/2)-2, (col / 2), username)
	stdscr.MovePrint(row/2, (col / 2), password)
	stdscr.MovePrint((row/2)+2, (col / 2), note)

	go func() {
		timer := time.NewTimer(time.Duration(sessionTimeout) * time.Second)
		for {
			select {
			case <-timer.C:
				stdscr.Erase()
				stdscr.MovePrint(row/2, col/2, "Decrypt Session Expired")
				stdscr.Refresh()
				return

			default:
				stdscr.MovePrint(1, 1, fmt.Sprintf("Shutting down in %d seconds...", sessionTimeout))
				sessionTimeout -= 1
				time.Sleep(1 * time.Second)
				stdscr.Refresh()
			}
		}
	}()
	stdscr.Refresh()

	stdscr.GetChar()

}

func getDefaultScreen() *gc.Window {
	stdscr, _ := gc.Init()

	if !gc.HasColors() {
		log.Fatal("Cryptorious requires a colour capable terminal")
	}

	// Must be called after Init but before using any colour related functions
	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(3, gc.C_BLUE, gc.C_BLACK)

	return stdscr
}
