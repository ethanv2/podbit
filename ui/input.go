package ui

import (
	"os"
	"unicode/utf8"

	"github.com/ethanv2/podbit/data"
	"github.com/ethanv2/podbit/sound"
)

var (
	exitChan chan int
)

func getInput(out chan rune, errc chan error) {
	var buf [4]byte
	var err error

	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		errc <- err
		out <- 0x0
	}

	c, _ := utf8.DecodeRune(buf[:])

	errc <- nil
	out <- c
}

// Exit requests that the input handler shuts down and gracefully
// exits the program via a return to the main function.
func Exit() {
	exitChan <- 1
}

// InputLoop - main UI input handler
//
// Receives all key inputs serially, one character at a time
// If there is no global keybinding for this key, we pass it
// to the UI subsystem, which can deal with it from there.
//
// Any and all key inputs causes an immediate and full UI redraw
func InputLoop(exit chan int) {
	exitChan = exit

	var c rune
	var char chan rune = make(chan rune)
	var err chan error = make(chan error, 1)

	for {
		go getInput(char, err)

		select {
		case c = <-char:
			if <-err != nil {
				return
			}

			switch c {

			case '1':
				ActivateMenu(PlayerMenu)
			case '2':
				ActivateMenu(QueueMenu)
			case '3':
				ActivateMenu(DownloadMenu)
			case '4':
				ActivateMenu(LibraryMenu)
			case 'r':
				data.Q.Reload()
			case 'p':
				sound.Plr.Toggle()
			case 's':
				sound.Plr.Stop()
			case 'c':
				sound.ClearQueue()
			case ']':
				sound.Plr.Seek(5)
			case '[':
				sound.Plr.Seek(-5)
			case '}':
				sound.Plr.Seek(60)
			case '{':
				sound.Plr.Seek(-60)
			case '\f': // Control-L
				UpdateDimensions(root, true)
			case 'q':
				if data.Caching.Ongoing() == 0 {
					return
				}

				StatusMessage("Error: Cannot quit with ongoing downloads")
			default:
				PassKeystroke(c)
			}

			Redraw(RedrawAll)
		case <-exit:
			return
		}
	}
}
