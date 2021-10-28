package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/ethanv2/podbit/colors"
	"github.com/ethanv2/podbit/data"
	"github.com/ethanv2/podbit/sound"

	"github.com/rthornton128/goncurses"
)

const (
	// MessageTime is the time a message will show for
	MessageTime time.Duration = 2 * time.Second
)

var (
	statusMessage = make(chan string)
	lastStatus time.Time

	status string
)

func trayWatcher() {
	for {
		time.Sleep(100 * time.Millisecond)
		Redraw(RedrawAll)
	}
}

// RenderTray renders the statusbar tray at the bottom of the screen
// Tray takes up two vertical cells and the entirety of the width
// The top cell is a horizontal line denoting a player status bar
// The bottom cell is the status text
func RenderTray(scr *goncurses.Window, w, h int) {
	scr.HLine(h-2, 0, goncurses.ACS_HLINE, w)

	pos, dur := sound.Plr.GetTimings()
	percent := pos / dur

	scr.ColorOn(colors.ColorBlue)
	scr.HLine(h-2, 0, '=', int(percent*float64(w)))

	head := int(math.Max((percent*float64(w))-1, 0))
	scr.MovePrint(h-2, head, ">")
	scr.ColorOff(colors.ColorBlue)

	now := time.Now()
	if now.Sub(lastStatus) > MessageTime {
		select {
		case status = <-statusMessage:
			lastStatus = now
		default:
			status = ""
		}
	}

	if status != "" {
		root.ColorOn(colors.ColorRed)
		scr.MovePrint(h-1, 0, status)
		root.ColorOff(colors.ColorRed)
	} else {
		pos, dur := sound.Plr.GetTimings()
		p, d := data.FormatTime(pos), data.FormatTime(dur)
		code := fmt.Sprintf("[%s/%s]", p, d)

		if sound.Plr.IsPlaying() {
			var status string
			if sound.Plr.IsPaused() {
				status = "Paused"
			} else {
				status = "Playing"
			}

			root.ColorOn(colors.ColorRed)
			scr.MovePrintf(h-1, 0, "%s: ", status)
			root.ColorOff(colors.ColorRed)

			root.ColorOn(colors.ColorBlue)
			scr.MovePrintf(h-1, len(status)+2, "%s - %s", sound.Plr.NowPlaying, sound.Plr.NowPodcast)
			root.ColorOff(colors.ColorBlue)

			root.ColorOn(colors.ColorRed)
			scr.MovePrintf(h-1, w-len(code), "%s", code)
			root.ColorOff(colors.ColorRed)
		} else if sound.Plr.IsWaiting() {
			root.ColorOn(colors.ColorRed)
			scr.MovePrint(h-1, 0, "Waiting for download...")
			root.ColorOff(colors.ColorRed)
		} else {
			root.ColorOn(colors.ColorRed)
			scr.MovePrint(h-1, 0, "Not playing")
			root.ColorOff(colors.ColorRed)
		}
	}
}

// StatusMessage sends a status message
//
// Will block for the previous message to finish first.
// Every message can be guaranteed MessageTime display time
func StatusMessage(msg string) {
	statusMessage <- msg
}
