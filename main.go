package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethanv2/podbit/colors"
	"github.com/ethanv2/podbit/data"
	"github.com/ethanv2/podbit/sound"
	"github.com/ethanv2/podbit/ui"

	"github.com/juju/fslock"
	"github.com/vit1251/go-ncursesw"
)

const (
	homebase = "podbit"
	pidfile  = "podbit.lock"
)

var (
	homedir string
	confdir string

	redraw    = make(chan int)
	keystroke = make(chan rune)
	newMen    = make(chan ui.Menu)
	exit      = make(chan int)
)

func banner() {
	fmt.Printf("Starting Podbit v%d.%d.%d...\n", verMaj, verMin, verPatch)
}

func initDirs() {
	share := os.Getenv("XDG_DATA_HOME")
	config, _ := os.UserConfigDir()

	if share == "" {
		share = ".local/share"
	}

	homedir = filepath.Join(share, homebase)
	confdir = filepath.Join(config, homebase)

	herr := os.MkdirAll(homedir, os.ModeDir|os.ModePerm)
	cerr := os.MkdirAll(confdir, os.ModeDir|os.ModePerm)

	if herr != nil || cerr != nil {
		fmt.Println("Error: Failed to required directory(s)")
		os.Exit(1)
	}
}

func alreadyRunning() (bool, *fslock.Lock) {
	lockpath := filepath.Join(homedir, pidfile)
	lock := fslock.New(lockpath)
	err := lock.TryLock()

	if err != nil {
		return true, nil
	}

	return false, lock
}

func initColors() {
	if goncurses.HasColors() {
		goncurses.StartColor()
	}

	colors.CreateColors()
}

func initTTY() {
	goncurses.Raw(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)
}

func main() {
	banner()
	initDirs()

	running, lock := alreadyRunning()
	if running {
		fmt.Println("Error: Podbit is already running")
		os.Exit(1)
	}
	defer lock.Unlock()

	now := time.Now()
	err := data.InitData()
	if err != nil {
		fmt.Println("\n" + err.Error())
		return
	}
	defer data.SaveData()
	go data.ReloadLoop()

	fmt.Print("Initialising sound system...")
	sound.Plr, err = sound.NewPlayer()
	if err != nil {
		fmt.Printf("\nError: Failed to initialise sound system: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("done")

	go sound.Mainloop()
	defer sound.Plr.Destroy()

	scr, err := goncurses.Init()
	if err != nil {
		fmt.Printf("Error: Failed to initialize UI: %s\n", err)
		os.Exit(1)
	}
	initTTY()
	initColors()
	defer goncurses.End()

	ui.InitUI(scr, ui.LibraryMenu, redraw, keystroke, newMen)
	go ui.RenderLoop()

	// Welcome message
	startup := time.Since(now)
	go ui.StatusMessage(fmt.Sprintf("Podbit v%d.%d.%d -- %d episodes of %d podcasts loaded in %.2fs",
		verMaj,
		verMin,
		verPatch,
		len(data.Q.Items), len(data.Q.GetPodcasts()),
		startup.Seconds()))

	// Initial UI draw
	ui.Redraw(ui.RedrawAll)

	// Initialisation is done; use this thread as the input loop
	ui.InputLoop(exit)
}
