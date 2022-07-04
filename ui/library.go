package ui

import (
	"fmt"

	"github.com/ethanv2/podbit/data"
	"github.com/ethanv2/podbit/sound"
	"github.com/ethanv2/podbit/ui/components"

	"github.com/vit1251/go-ncursesw"
)

// Library represents the list menu type and state.
//
// Library displays all detected and configured podcasts, along
// with associated episodes sorted into said podcasts.
type Library struct {
	men [2]components.Menu

	menSel int
}

func (l *Library) Name() string {
	return "Library"
}

func (l *Library) renderPodcasts(x, y int) {
	l.men[0].X = x
	l.men[0].Y = y

	l.men[0].W, l.men[0].H = (w/2)-1, (h - 5)
	l.men[0].Win = *root

	l.men[0].Items = l.men[0].Items[:0]
	l.men[0].Items = data.Q.GetPodcasts()

	l.men[0].Selected = true

	if len(l.men[0].Items) > 0 {
		l.men[0].Render()
	} else {
		root.MovePrint(y, x, "No podcasts")
	}
}

func (l *Library) renderEpisodes(x, y int) {
	if len(l.men[0].Items) < 1 {
		return
	}

	l.men[1].X = x
	l.men[1].Y = y

	l.men[1].W, l.men[1].H = (w/2)-2, (h - 5)
	l.men[1].Win = *root

	l.men[1].Items = l.men[1].Items[:0]

	data.Q.RevRange(func(i int, elem *data.QueueItem) bool {
		_, sel := l.men[0].GetSelection()
		if data.DB.GetFriendlyName(elem.URL) == sel {
			var text string
			entry, ok := data.Downloads.Query(elem.Path)
			title := entry.Title
			if !ok || title == "" {
				text = elem.URL
			} else {
				text = title
			}

			l.men[1].Items = append(l.men[1].Items, text)
		}

		return true
	})

	l.men[1].Selected = (l.menSel == 1)

	l.men[1].Render()
}

func (l *Library) Render(x, y int) {
	l.renderPodcasts(x, y)

	root.AttrOn(goncurses.A_BOLD)
	root.VLine(y, w/2, goncurses.ACS_VLINE, h-2-y)
	root.AttrOff(goncurses.A_BOLD)

	l.renderEpisodes(w/2+1, y)
}

func (l *Library) Input(c rune) {
	switch c {
	case 'j':
		l.men[l.menSel].MoveSelection(1)
	case 'k':
		l.men[l.menSel].MoveSelection(-1)
	case 'h':
		l.MoveSelection(-1)
	case 'l':
		l.MoveSelection(1)
	case 'g':
		l.men[l.menSel].ChangeSelection(0)
	case 'G':
		l.men[l.menSel].ChangeSelection(len(l.men[l.menSel].Items) - 1)
	case ' ':
		l.StartDownload()
	case 13:
		l.StartPlaying(false) // Enter key - enqueue
	case '\t':
		l.StartPlaying(true) // Tab key - play NOW!
	}
}

func (l *Library) ChangeSelection(index int) {
	if index >= len(l.men) || index < 0 {
		return
	}

	l.menSel = index

	if l.menSel == 0 {
		l.men[1].ChangeSelection(0)
	}
}

func (l *Library) MoveSelection(direction int) {
	if direction == 0 {
		return
	}

	off := l.menSel + direction
	l.ChangeSelection(off)
}

// StartDownload downloads the currently focused library entry.
func (l *Library) StartDownload() {
	if len(l.men[0].Items) < 1 || len(l.men[1].Items) < 1 {
		return
	}

	defer func() {
		// Move cursor down
		l.men[l.menSel].MoveSelection(1)
	}()

	targets := l.men[1].Items
	if l.menSel == 1 {
		_, target := l.men[1].GetSelection()
		item := data.Q.GetEpisodeByURL(target)

		if item == nil {
			return
		}

		if y, _ := data.Downloads.IsDownloading(item.Path); y {
			go StatusMessage("Episode already downloading")
			return
		}

		data.Downloads.Download(item)
		go StatusMessage(fmt.Sprintf("Download of %s started...", item.URL))

		return
	}

	for _, elem := range targets {
		if data.IsURL(elem) {
			item := data.Q.GetEpisodeByURL(elem)
			if item == nil {
				continue
			}

			if y, _ := data.Downloads.IsDownloading(item.Path); y {
				go StatusMessage("Episode already downloading")
				return
			}

			go data.Downloads.Download(item)
		}
	}

	go StatusMessage("Download of multiple episodes started...")
}

// StartPlaying begins playing the currently focused element.
// If the current focus requires downloading (and enough information
// is known to oblige) it will first be downloaded.
func (l *Library) StartPlaying(immediate bool) {
	if len(l.men[0].Items) < 1 || len(l.men[1].Items) < 1 {
		return
	}

	defer func() {
		// Move cursor down
		l.men[l.menSel].MoveSelection(1)
	}()

	if l.menSel == 1 {
		var item *data.QueueItem
		_, entry := l.men[1].GetSelection()
		if data.IsURL(entry) {
			item = data.Q.GetEpisodeByURL(entry)
		} else {
			item = data.Q.GetEpisodeByTitle(entry)
		}

		if immediate {
			go StatusMessage(fmt.Sprintf("Now playing episode %q", entry))
			sound.PlayNow(item)
		} else {
			go StatusMessage(fmt.Sprintf("Enqueued episode %q to play", entry))
			sound.Enqueue(item)
		}

	} else {
		if immediate {
			return
		}

		_, entry := l.men[0].GetSelection()
		sound.EnqueueByPodcast(entry)
		go StatusMessage("Multiple episodes enqueued...")
	}
}
