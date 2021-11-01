package ui

import (
	"fmt"
	"strconv"

	"github.com/ethanv2/podbit/colors"
	"github.com/ethanv2/podbit/data"
	"github.com/ethanv2/podbit/sound"
	"github.com/ethanv2/podbit/ui/components"
)

var downloadHeadings []components.Column = []components.Column{
	{
		Label: "ID",
		Width: 0.1,
		Color: colors.BackgroundGreen,
	},
	{
		Label: "%",
		Width: 0.1,
		Color: colors.BackgroundYellow,
	},
	{
		Label: "Path",
		Width: 0.4,
		Color: colors.BackgroundBlue,
	},
	{
		Label: "Status",
		Width: 0.4,
		Color: colors.BackgroundRed,
	},
}

type Downloads struct {
	showDone bool

	tbl components.Table
}

func (q *Downloads) Name() string {
	return "Downloads"
}

func (q *Downloads) Render(x, y int) {
	q.tbl.X, q.tbl.Y = x, y
	q.tbl.W, q.tbl.H = w, h
	q.tbl.Win = root

	q.tbl.Columns = downloadHeadings

	q.tbl.Items = nil
	for i, elem := range data.Downloads.Downloads() {
		item := make([]string, len(downloadHeadings))

		item[0] = strconv.FormatInt(int64(i), 10)
		item[1] = strconv.FormatFloat(elem.Percentage*100, 'f', 2, 64)
		item[2] = elem.Path

		if elem.Completed {
			if elem.Success {
				item[3] = "Finished"
			} else {
				item[3] = fmt.Sprintf("Failed (%s)", elem.Error)
			}
		} else {
			item[3] = "In progress"
		}

		q.tbl.Items = append(q.tbl.Items, item)
	}

	q.tbl.Render()
}

func (q *Downloads) Input(c rune) {
	switch c {
	case 'j':
		q.tbl.MoveSelection(1)
	case 'k':
		q.tbl.MoveSelection(-1)
	case 'd':
		q.Cancel()
	case 13:
		q.Enqueue()
	}
}

func (q *Downloads) Enqueue() {
	i, _ := q.tbl.GetSelection()
	d := data.Downloads.Downloads()[i].Path

	var found *data.QueueItem
	data.Q.Range(func(i int, item *data.QueueItem) bool {
		if item.Path == d {
			found = item
		}

		return true
	})

	if found != nil {
		go StatusMessage("Enqueued: Download will play once completed")
		sound.Enqueue(found)
	}
}

func (q *Downloads) Cancel() {
	i, _ := q.tbl.GetSelection()
	if i >= len(data.Downloads.Downloads()) {
		return
	}

	dl := data.Downloads.Downloads()[i]
	if !dl.Completed {
		dl.Stop <- 1
	}
}
