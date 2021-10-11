// Package sound is responsible for playing audio and managing mpv instances
//
// This package is usually interacted with by the front-end UI code to add
// items to the queue. The rest happens automatically using multi-threaded
// player logic. All that needs to be maintained from the interface side is
// the queue and the rest is run automatically.
//
// The queue is a simple FIFO structure formed from a slice of QueueEntries.
// Entries that require downloading are handled gracefully and with as little
// user impact as possible. Usually, the user won't even notice anything happened.
//
// Sound is played through an idle MPV instance which sits in the background and
// recieves media to play when appropriate
package sound

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/ethanv2/podbit/data"

	"github.com/blang/mpv"
)

// Useful player vars
var (
	// PlayerName is the name of the player program to spawn
	PlayerName = "mpv"
	// The path to the RPC endpoint
	PlayerRPC = "/tmp/podbit-mpv"
	// PlayerArgs are the standard arguments to use for the player
	// These are not the final configs of the player, but just used
	// to idle mpv ready to recieve instructions
	PlayerArgs = []string{"--no-video", "--idle", "--input-ipc-server=" + PlayerRPC}
	// UpdateTime is the time between queue checks and supervision updates
	UpdateTime = 200 * time.Millisecond
)

// Player represents the current player instance
type Player struct {
	proc *exec.Cmd

	exit      chan int
	stop      chan int
	watchStop chan int

	ipcc *mpv.IPCClient
	ctrl *mpv.Client

	output io.ReadCloser
	times  io.ReadCloser

	Waiting  bool
	download *data.Download

	Playing bool

	NowPlaying string
	NowPodcast string
}

var (
	Plr Player
)

func updateWait(u chan int) {
	time.Sleep(UpdateTime)
	u <- 1
}

// Detect end of process and exit if it does
func pin(p *Player, giveUp chan int) {
	p.proc.Wait()

	// Uh oh! We exitted mpv - now we need to exit quickly too
	p.exit <- 1
}

func NewPlayer(exit chan int) (p Player, err error) {
	p.exit = exit

	p.proc = exec.Command(PlayerName, PlayerArgs...)
	p.output, err = p.proc.StdoutPipe()
	p.times, err = p.proc.StderrPipe()
	p.proc.Start()

	for err = ConnectPlayer(&p); err != nil; {
		err = ConnectPlayer(&p)
	}

	p.watchStop = make(chan int)
	p.stop = make(chan int)
	go pin(&p, p.watchStop)

	return
}

// ConnectPlayer attempts to connect to the RPC endpoint
// Sadly, this is needed because of an exceptionally bad
// design choice in the mpv library forcing me to create
// this bad workaround. :(
func ConnectPlayer(p *Player) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error: Player connection")
		}
	}()

	p.ipcc = mpv.NewIPCClient(PlayerRPC)
	p.ctrl = mpv.NewClient(p.ipcc)

	return
}

func (p *Player) Play(q *data.QueueItem) {
	if q.State != data.StatePending {
		now, ok := data.Caching.Query(q.Path)
		if !ok {
			p.NowPlaying = ""
			p.NowPodcast = ""
		}
		p.NowPlaying = now.Title
		p.NowPodcast = data.DB.GetFriendlyName(q.URL)

		p.ctrl.Loadfile(q.Path, mpv.LoadFileModeReplace)
		p.Playing = true
	}
}

func (p *Player) Stop() {
	p.ctrl.SetPause(true)
	p.ctrl.Seek(0, mpv.SeekModeAbsolute)
	p.Playing = false
}

func (p *Player) Destroy() {
	p.proc.Process.Kill()
	p.Playing = false
}

func (p *Player) IsPaused() bool {
	paused, _ := p.ctrl.Pause()
	return paused
}

func (p *Player) Pause() {
	if !p.Playing {
		return
	}

	// Leave playing set to true so we know not to play another episode
	p.ctrl.SetPause(true)
}

func (p *Player) Unpause() {
	if !p.Playing {
		return
	}

	// Leave playing set to true so we know not to play another episode
	p.ctrl.SetPause(false)
}

func (p *Player) Toggle() {
	paused, _ := p.ctrl.Pause()
	p.ctrl.SetPause(!paused)
}

// GetTimings returns the current time and duration
// of the ongoing player. Returns zero if we are
// not playing currently
func (p *Player) GetTimings() (float64, float64) {
	if !p.Playing {
		return 0, 0
	}

	pos, _ := p.ctrl.Position()
	dur, _ := p.ctrl.Duration()

	return pos, dur
}

// Wait for the current episode to complete
func (p *Player) Wait() {
	if !p.Playing {
		return
	}

	// Wait for media to load?
	var dur float64
	for dur == 0 {
		_, dur = p.GetTimings()
	}

	time.Sleep(time.Duration(dur+1) * time.Second)
}

func Mainloop() {
	for {
		if !Plr.Playing && !Plr.Waiting && len(queue) > 0 {
			elem, stop := GetQueueHead()
			if stop {
				Plr.Stop()
				continue
			}

			if elem.State != data.StatePending && data.Caching.EntryExists(elem.Path) {
				Plr.Play(elem)
				Plr.Wait()

				Plr.Playing = false
			} else {
				Plr.Waiting = true

				id, _ := data.Caching.Download(elem)
				Plr.download = &data.Caching.Downloads[id]
				for !Plr.download.Completed {
				}

				Plr.Waiting = false
				head--
				continue
			}
		}

		u := make(chan int)
		go updateWait(u)
		select {
		case <-u:
		case <-Plr.stop:
			return
		}
	}
}
