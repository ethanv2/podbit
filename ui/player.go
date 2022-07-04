package ui

import (
	"math"

	"github.com/ethanv2/podbit/colors"
	"github.com/ethanv2/podbit/data"
	"github.com/ethanv2/podbit/sound"

	"github.com/vit1251/go-ncursesw"
)

// Player is the full screen player component.
//
// Player displays the currently playing episode, the next up
// episode, progress through the episode etc.
//
// This is mostly for user convenience and visual appeal.
type Player struct{}

func (l *Player) Name() string {
	return "Now playing"
}

func (l *Player) Render(x, y int) {
	pos, dur := sound.Plr.GetTimings()
	percent := pos / dur
	p, d := data.FormatTime(pos), data.FormatTime(dur)
	ep, pod := data.LimitString(sound.Plr.NowPlaying, w-1), data.LimitString(sound.Plr.NowPodcast, w-1)

	minxt := int(math.Max(0, float64((w-len([]rune(ep)))/2)))
	minxp := int(math.Max(0, float64((w-len([]rune(pod)))/2)))

	var stat, div string
	if sound.Plr.IsPlaying() {
		stat = "Now Playing: "
		div = "by"
	} else {
		stat = "Not playing"
		div = ""
	}

	// Now playing
	root.ColorOn(colors.ColorRed)
	root.MovePrint(x+4, (w-len(stat))/2, stat)
	root.ColorOff(colors.ColorRed)

	// [episode name]
	root.AttrOn(goncurses.A_BOLD)
	root.MovePrint(x+6, minxt, ep)
	root.AttrOff(goncurses.A_BOLD)

	// by
	root.ColorOn(colors.ColorGreen)
	root.MovePrint(x+8, (w-len(div))/2, div)
	root.ColorOff(colors.ColorGreen)

	// [podcast name]
	root.ColorOn(colors.ColorYellow)
	root.MovePrint(x+10, minxp, pod)
	root.ColorOff(colors.ColorGreen)

	root.MovePrint(h-(h/3), 0, p)
	root.MovePrint(h-(h/3), w-len(d), d)

	root.HLine(h-(h/3), len(p)+1, goncurses.ACS_HLINE, w-len(d)-len(p)-2)
	root.MovePrint(h-(h/3), len(p), "|")
	root.MovePrint(h-(h/3), len(p)+1+(w-len(d)-len(p)-2), "|")

	wid := int(float64(w-len(d)-len(p)-2) * percent)

	root.ColorOn(colors.ColorBlue)
	root.HLine(h-(h/3), len(p)+1, goncurses.ACS_HLINE, wid)
	root.ColorOff(colors.ColorBlue)

	// Next up
	cur, _ := sound.GetNext()
	lbl := "Next up: "
	if cur != nil {
		var txt string
		if entry, ok := data.Downloads.Query(cur.Path); ok && entry.Title != "" {
			txt = entry.Title
		} else {
			txt = cur.URL
		}
		clipped := data.LimitString(txt, w-1)

		root.ColorOn(colors.ColorRed)
		root.MovePrint(h-(h/3)+2, (w-len(lbl))/2, lbl)
		root.ColorOff(colors.ColorRed)

		root.MovePrint(h-(h/3)+3, (w-len(clipped))/2, clipped)
	}
}

func (l *Player) Input(c rune) {
	switch c {
	case ' ':
		sound.Plr.Toggle()
	case 'h':
		sound.Plr.Seek(-5)
	case 'l':
		sound.Plr.Seek(5)

	case 'H':
		sound.Plr.Seek(-60)
	case 'L':
		sound.Plr.Seek(60)
	}
}
