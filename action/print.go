package action

import (
	"fmt"
	"log"

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

	stdscr.Refresh()
	stdscr.GetChar()
}
