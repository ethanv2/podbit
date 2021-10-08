package data

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Possible cache errors
var (
	CacheIOError        error  = errors.New("Error: Failed to create cache entry")
	CacheDownloadFailed string = "Error: Failed to download from url %s"
)

type Cache struct {
	dir string

	episodes sync.Map

	downloadsMutex sync.RWMutex
	Downloads      []Download
}

type Episode struct {
	entry *QueueItem

	title string
}

type Download struct {
	episode Episode
	file    *os.File

	percentage float64

	size int64
	done int64

	started time.Time

	completed bool
	success   bool
}

func (c *Cache) progressWatcher(watch *Download, stop chan int) {
	for {
		c.downloadsMutex.Lock()

		fi, err := watch.file.Stat()
		if err != nil {
			return
		}
		watch.done = fi.Size()
		watch.percentage = float64(watch.done) / float64(watch.size)

		c.downloadsMutex.Unlock()

		select {
		case <-stop:
			return
		default:
			time.Sleep(1 * time.Second)

			fmt.Printf("%f%% (%d/%d) - elapsed: %s\n", watch.percentage, watch.done, watch.size, time.Now().Sub(watch.started))
		}
	}
}

// Dig through newsboat stuff to guess the download dir
// If we can't find it, just use the newsboat default and hope for the best
func (c *Cache) guessDir() string {
	conf, _ := os.UserConfigDir()
	p := filepath.Join(conf, "newsboat/config")

	file, err := os.Open(p)
	if err != nil {
		ret, _ := os.UserHomeDir()
		return ret
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			ret, _ := os.UserHomeDir()
			return ret
		}

		line := scanner.Text()
		fields := strings.Split(line, " ")

		if len(line) < 1 || len(fields) < 2 {
			continue
		}

		if fields[0] == "download-path" {
			return fields[1]
		}
	}

	ret, _ := os.UserHomeDir()
	return ret
}

// Open and initialise the cache
func (c *Cache) Open() error {
	c.dir = c.guessDir()

	return nil
}

// Start a download and return its ID in the downloads table
// This can be used to retrieve information about said download
//
// WARNING: DO NOT modify this table without taking out the mutex
// This code is NOT THREAD SAFE witout this mutex being used
//
// Returns as soon as the download has been initialised
// Does not block until completion
func (c *Cache) Download(item *QueueItem) (id int, err error) {

	f, err := os.Create(item.Path)
	if err != nil {
		return 0, CacheIOError
	}
	defer f.Close()

	resp, err := http.Get(item.Url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(CacheDownloadFailed, item.Url)
	}
	defer resp.Body.Close()

	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	newEp := Episode{
		entry: item,
	}

	stop := make(chan int)

	c.downloadsMutex.Lock()
	var dl Download = Download{
		episode: newEp,
		file:    f,
		size:    size,
		started: time.Now(),
	}
	c.Downloads = append(c.Downloads, dl)

	id = len(c.Downloads) - 1
	go c.progressWatcher(&c.Downloads[id], stop)
	c.downloadsMutex.Unlock()

	defer func() {
		stop <- 1
	}()

	go func() {
		io.Copy(f, resp.Body)

		c.downloadsMutex.Lock()
		c.Downloads[id].completed = true
		if err != nil {
			c.Downloads[id].success = false
		} else {
			c.Downloads[id].success = true
		}

		c.episodes.Store(item.Path, newEp)
	}()

	return
}
