package colors

import (
	"github.com/vit1251/go-ncursesw"
)

// Color pair ID definitions.
const (
	ColorRed = iota + 1
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan

	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
)

// CreateColors initialises all colors for ncurses usage.
// Should be called only once per application run.
func CreateColors() {
	goncurses.InitPair(ColorRed, goncurses.C_RED, goncurses.C_BLACK)
	goncurses.InitPair(ColorGreen, goncurses.C_GREEN, goncurses.C_BLACK)
	goncurses.InitPair(ColorYellow, goncurses.C_YELLOW, goncurses.C_BLACK)
	goncurses.InitPair(ColorBlue, goncurses.C_BLUE, goncurses.C_BLACK)
	goncurses.InitPair(ColorMagenta, goncurses.C_MAGENTA, goncurses.C_BLACK)
	goncurses.InitPair(ColorCyan, goncurses.C_CYAN, goncurses.C_BLACK)

	goncurses.InitPair(BackgroundRed, goncurses.C_BLACK, goncurses.C_RED)
	goncurses.InitPair(BackgroundGreen, goncurses.C_BLACK, goncurses.C_GREEN)
	goncurses.InitPair(BackgroundYellow, goncurses.C_BLACK, goncurses.C_YELLOW)
	goncurses.InitPair(BackgroundBlue, goncurses.C_BLACK, goncurses.C_BLUE)
	goncurses.InitPair(BackgroundMagenta, goncurses.C_BLACK, goncurses.C_MAGENTA)
	goncurses.InitPair(BackgroundCyan, goncurses.C_BLACK, goncurses.C_CYAN)
}
