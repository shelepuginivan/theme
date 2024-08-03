// Package theme provides theme management capabilities.
package theme

import (
	"encoding/json"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adrg/xdg"
	cp "github.com/otiai10/copy"
	"github.com/shelepuginivan/theme/errorlist"
)

const (
	// File containing [Copy] entries.
	copyfile = "copy.json"

	// Directory that contains executable files.
	// These executables are ran during the theme setting.
	rundir = "run.d"
)

// Copy represents an entry that should be copied.
type Copy struct {
	Src string `json:"src"` // What to copy (source).
	Dst string `json:"dst"` // Where to copy (destination).
}

// Themer manages themes.
type Themer struct {
	prefix string // Prefix directory where themes are stored.
	el     *errorlist.Errorlist
}

// New returns a new instance of [Themer].
func New() *Themer {
	return NewWithPrefix(filepath.Join(xdg.ConfigHome, "theme"))
}

// NewWithPrefix is like [New] but allows to specify prefix directory.
func NewWithPrefix(prefix string) *Themer {
	return &Themer{
		prefix: prefix,
		el:     errorlist.New(),
	}
}

// Set sets a theme by name. It copies files from `copy.json` and runs `run`
// located in the theme directory.
// Returned slice is a slice of errors occurred during copying and execution.
func (t *Themer) Set(name string) []error {
	c, err := t.ReadCopy(name)
	if err != nil {
		t.el.Append(err)
	}

	var wg sync.WaitGroup

	for _, e := range c {
		wg.Add(1)
		go func() {
			err := cp.Copy(t.ExpandPath(e.Src, name), t.ExpandPath(e.Dst, name))
			if err != nil {
				t.el.Append(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	entries, err := t.ReadRun(name)
	if err != nil {
		t.el.Append(err)
	}

	for _, e := range entries {
		wg.Add(1)
		go func() {
			cmd := exec.Command(e)
			cmd.Dir = filepath.Join(t.prefix, name)

			if err := cmd.Run(); err != nil {
				t.el.Append(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return t.el.Get()
}

// Themes returns a slice of available themes.
func (t *Themer) Themes() (themes []string, err error) {
	entries, err := os.ReadDir(t.prefix)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			themes = append(themes, e.Name())
		}
	}

	return themes, nil
}

// Random returns a random theme.
func (t *Themer) Random() (string, error) {
	themes, err := t.Themes()
	if err != nil {
		return "", err
	}

	return themes[rand.Intn(len(themes))], nil
}

// ReadCopy reads `copy.json` located in the theme directory and returns
// [ReadCopy] entries.
func (t *Themer) ReadCopy(name string) (c []Copy, err error) {
	data, err := os.ReadFile(filepath.Join(t.prefix, name, copyfile))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// ReadRun reads `run.d` located in the theme directory and returns paths to
// all files in it.
func (t *Themer) ReadRun(name string) (r []string, err error) {
	rundirPath := filepath.Join(t.prefix, name, rundir)

	entries, err := os.ReadDir(rundirPath)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		r = append(r, filepath.Join(rundirPath, e.Name()))
	}
	return r, nil
}

// ExpandPath expands path for convenience.
//   - Environment variables are expanded
//   - `~` is replaced with user home directory
//   - `@` is replaced with theme directory
func (t *Themer) ExpandPath(path, name string) string {
	r := strings.NewReplacer(
		"@", filepath.Join(t.prefix, name),
		"~", xdg.Home,
	)

	return os.ExpandEnv(r.Replace(path))
}
