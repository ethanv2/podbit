package components

import (
	"github.com/ethanv2/podbit/colors"
	"github.com/vit1251/go-ncursesw"
)

// Menu represents a vertical panel menu of cellular entries
// taking up an entire row from X to W
//
// Menu handles focus, scrolling and the managing of elements
// at each render.
type Menu struct {
	X, Y int
	Win  goncurses.Window

	Selected     bool
	prevw, prevh int

	List[string]
}

// Render immediately renders the menu to the specified
// fields X, Y, H and W.
func (m *Menu) Render() {
	// Reset scroll if invalid or not selected
	if len(m.Items) <= m.scroll || !m.Selected {
		m.scroll = 0
		m.sel = 0
	}

	items := m.Items[m.scroll:]
	c := m.scroll

	if m.prevw != m.W || m.prevh != m.H {
		m.scroll = 0
	}
	m.prevw, m.prevh = m.W, m.H

	if m.scroll+m.H < m.sel {
		m.sel = m.H
	}

	for i, elem := range items {
		if c == m.sel && m.Selected {
			m.Win.ColorOn(colors.BackgroundBlue)
		} else if i > m.H {
			break
		}

		capped, decode := elem, []rune(elem)
		if len(decode) > m.W {
			capped = string(decode[len(decode)-m.W:])
			capped = "<" + capped
		} else {
			// Pad out to fill row
			for i := len(decode); i <= m.W; i++ {
				capped += " "
			}
		}

		m.Win.MovePrint(m.Y+i, m.X, capped)
		c++

		m.Win.ColorOff(colors.BackgroundBlue)
	}
}
