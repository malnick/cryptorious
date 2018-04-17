package action

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/cryptorious/config"
	"github.com/malnick/cryptorious/vault"
	gc "github.com/rthornton128/goncurses"
)

// PrintAll will print the entire vault contents using ncurses menu
func PrintAll(c config.Config) error {
	log.Info("Opening session for all entries in vault")
	cleartextValues, err := decryptVault(c)
	if err != nil {
		return err
	}

	menuItems := []string{}
	for key := range cleartextValues {
		menuItems = append(menuItems, key)
	}

	printWithMenu(menuItems, cleartextValues, c.DecryptSessionTimeout)
	return nil
}

func decryptVault(c config.Config) (map[string]vault.CleartextEntry, error) {
	cleartextEntries := map[string]vault.CleartextEntry{}

	v, err := vault.New(c.VaultPath)
	if err != nil {
		return cleartextEntries, err
	}

	for k, e := range v.Data {
		clr, err := e.Decrypt(c.KMSClient)
		if err != nil {
			return cleartextEntries, err
		}

		cleartextEntries[k] = clr
	}

	return cleartextEntries, nil
}

func printWithMenu(menuItems []string, cleartextEntries map[string]vault.CleartextEntry, timeout int) error {
	defer gc.End()
	stdscr := getDefaultScreen()

	const (
		HEIGHT = 10
		WIDTH  = 40
	)

	gc.StartColor()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

	items := make([]*gc.MenuItem, len(menuItems))
	for i, val := range menuItems {
		items[i], _ = gc.NewItem(val, "")
		defer items[i].Free()
	}

	// create the menu
	menu, _ := gc.NewMenu(items)
	defer menu.Free()

	menuwin, _ := gc.NewWindow(HEIGHT, WIDTH, 4, 14)
	menuwin.Keypad(true)

	menu.SetWindow(menuwin)
	dwin := menuwin.Derived(6, 38, 3, 1)
	menu.SubWindow(dwin)
	//menu.Option(gc.O_SHOWDESC, true)
	menu.Format(5, 1)
	menu.Mark(" * ")

	// MovePrint centered menu title
	title := "Cryptorious Vault"
	menuwin.Box(0, 0)
	menuwin.ColorOn(1)
	menuwin.MovePrint(1, (WIDTH/2)-(len(title)/2), title)
	menuwin.ColorOff(1)
	menuwin.MoveAddChar(2, 0, gc.ACS_LTEE)
	menuwin.HLine(2, 1, gc.ACS_HLINE, WIDTH-2)
	menuwin.MoveAddChar(2, WIDTH-1, gc.ACS_RTEE)

	y, _ := stdscr.MaxYX()
	stdscr.ColorOn(2)
	stdscr.MovePrint(y-3, 1,
		"Use up/down arrows or page up/down to navigate. 'q' to exit")
	stdscr.ColorOff(2)
	stdscr.Refresh()

	menu.Post()
	defer menu.UnPost()
	menuwin.Refresh()

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
			password := cleartextEntries[entry].Password
			username := cleartextEntries[entry].Username
			note := cleartextEntries[entry].SecureNote
			printDecrypted(entry, username, password, note, timeout)
			return nil
		default:
			menu.Driver(gc.DriverActions[ch])
			menuwin.Refresh()
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
	for i := range sep {
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
				sessionTimeout--
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
