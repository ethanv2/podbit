package data

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dhowden/tag"
)

// Possible cache errors.
var (
	ErrorIO             = errors.New("Error: Failed to create cache entry")
	ErrorCreation       = errors.New("Error: Download directory did not exist and could not be created")
	ErrorDownloadFailed = "Error: Failed to download from url %s"
)

// Cache is the current state of the on-disk cache and associated
// operations.
//
// This structure *is* thread safe, but ONLY is used with the
// correct methods. Use with care!
type Cache struct {
	dir string

	episodes sync.Map

	downloadsMutex sync.RWMutex // Protects the below two variables
	downloads      []*Download
	ongoing        int
}

// Episode represents the data extracted from a single cached episode
// media entry.
type Episode struct {
	Queued bool

	Title string
	Date  int
	Host  string
}

// Dig through newsboat stuff to guess the download dir.
// If we can't find it, just use the newsboat default and hope for the best.
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

// Open opens and initialises the cache.
// Should be called once and once only - further modifications
// and cache mutations happen exclusively through other methods.
func (c *Cache) Open() error {
	home, _ := os.UserHomeDir()
	c.dir = strings.ReplaceAll(c.guessDir(), "~", home)
	files, err := ioutil.ReadDir(c.dir)

	if err != nil {
		cerr := os.MkdirAll(c.dir, os.ModeDir|os.ModePerm)
		if cerr != nil {
			return ErrorCreation
		}
	}

	for _, elem := range files {
		path := filepath.Join(c.dir, elem.Name())
		c.loadFile(path, true)
	}

	return nil
}

func (c *Cache) loadFile(path string, startup bool) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	defer func() {
		// Prevent invalid media files from causing a panic
		if p := recover(); p != nil {
			fmt.Printf("\nInvalid media file %q in cache! Ignoring...\n", path)
			return
		}
	}()

	data, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("\nError: Invalid media file %q in cache! Ignoring...\n", path)
		return
	}

	artist, albumArtist := data.Artist(), data.AlbumArtist()
	var host string
	if artist == "" {
		host = albumArtist
	} else {
		host = artist
	}

	ep := Episode{
		Queued: !startup,
		Title:  data.Title(),
		Date:   data.Year(),
		Host:   host,
	}

	c.episodes.Store(path, ep)
}

// Download starts an asynchronous download in a new goroutine.
// Returns the ID in the downloads table, which must be accessed
// using a mutex.
func (c *Cache) Download(item *QueueItem) (id int, err error) {
	f, err := os.Create(item.Path)
	dl := Download{
		mut:     new(sync.RWMutex),
		Path:    item.Path,
		File:    f,
		Elem:    item,
		Started: time.Now(),
		Stop:    make(chan int),
	}

	if err != nil {
		dl = Download{
			Started:   time.Now(),
			Completed: true,
			Success:   false,
			Elem:      item,
			Error:     "IO Error",
		}

		c.downloadsMutex.Lock()
		c.downloads = append(c.downloads, &dl)
		id = len(c.downloads) - 1
		c.downloadsMutex.Unlock()

		return id, ErrorIO
	}

	if item.Youtube {
		go dl.DownloadYoutube()
	} else {
		go dl.DownloadHTTP()
	}

	c.downloadsMutex.Lock()
	c.downloads = append(c.downloads, &dl)
	id = len(c.downloads) - 1
	c.ongoing++
	c.downloadsMutex.Unlock()

	return
}

// IsDownloading queries the download cache to check.
// if a podcast is currently downloading.
func (c *Cache) IsDownloading(path string) (bool, int) {
	c.downloadsMutex.RLock()
	defer c.downloadsMutex.RUnlock()

	for i, elem := range c.downloads {
		elem.mut.RLock()
		if elem.Path == path && !elem.Completed {
			elem.mut.RUnlock()
			return true, i
		}
		elem.mut.RUnlock()
	}

	return false, 0
}

// GetDownload returns the specified download in a thread-safely.
// This should be used to get the details of a specified download
// via the ID.
func (c *Cache) GetDownload(ind int) Download {
	c.downloads[ind].mut.RLock()
	defer c.downloads[ind].mut.RUnlock()

	return *c.downloads[ind]
}

// Ongoing returns the current number of ongoing downloads.
// The value cannot change while this function is executing.
func (c *Cache) Ongoing() int {
	c.downloadsMutex.RLock()
	defer c.downloadsMutex.RUnlock()

	return c.ongoing
}

// Downloads returns all recorded downloads at this point,
// including completed or failed downloads.
// No downloads can start or end while this function is
// executing.
func (c *Cache) Downloads() []Download {
	c.downloadsMutex.RLock()
	defer c.downloadsMutex.RUnlock()

	dls := make([]Download, len(c.downloads))

	for i, elem := range c.downloads {
		elem.mut.RLock()
		dls[i] = *elem
		elem.mut.RUnlock()
	}

	return dls
}

// Query returns cached data about an episode on disk.
func (c *Cache) Query(path string) (ep Episode, ok bool) {
	e, ok := c.episodes.Load(path)
	if e != nil {
		ep = e.(Episode)
	}

	return
}

// QueryAll returns all known data about the on-disk cache.
func (c *Cache) QueryAll(allowQueued bool) (e []Episode) {
	c.episodes.Range(func(key interface{}, value interface{}) bool {
		ep := value.(Episode)
		if (!allowQueued && !ep.Queued) || allowQueued {
			e = append(e, ep)
		}

		return true
	})

	return
}

// EntryExists searches the cache to determine if the entry exists
// Path should be an absolute path
// If path lies outside the cache dir, false is returned.
func (c *Cache) EntryExists(path string) bool {
	f, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	f.Close()
	return true
}
