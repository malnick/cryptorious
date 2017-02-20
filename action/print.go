package action

import (
	"fmt"
	"log"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

func printDecrypted(key, username, password, note string) {
	stdscr, _ := gc.Init()
	defer gc.End()

	if !gc.HasColors() {
		log.Fatal("Example requires a colour capable terminal")
	}

	// Must be called after Init but before using any colour related functions
	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	stdscr.ColorOn(1)

	stdscr.Border(gc.ACS_VLINE, gc.ACS_VLINE, gc.ACS_HLINE, gc.ACS_HLINE,
		gc.ACS_ULCORNER, gc.ACS_URCORNER, gc.ACS_LLCORNER, gc.ACS_LRCORNER)

	row, col := stdscr.MaxYX()
	keyMsg := fmt.Sprintf("Decrypted Vault Entry for %s", key)
	passwordMsg := fmt.Sprintf("Password: %s", password)
	usernameMsg := fmt.Sprintf("Username: %s", username)
	noteMsg := fmt.Sprintf("Secure Note: %s", note)

	sep := make([]string, len(keyMsg))
	for i, _ := range sep {
		sep[i] = "-"
	}

	stdscr.MovePrint((row/2)-3, (col/2)-len(keyMsg), keyMsg)
	stdscr.MovePrint((row/2)-2, (col/2)-len(keyMsg), strings.Join(sep, ""))
	stdscr.MovePrint((row/2)-1, (col/2)-len(usernameMsg), usernameMsg)
	stdscr.MovePrint(row/2, (col/2)-len(passwordMsg), passwordMsg)
	stdscr.MovePrint((row/2)+1, (col/2)-len(noteMsg), noteMsg)

	stdscr.Refresh()
	stdscr.GetChar()
}
