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
)

const (
	// File containing [Copy] entries.
	copyfile = "copy.json"

	// File that is ran during the theme setting.
	runfile = "run"
)

// Copy represents an entry that should be copied.
type Copy struct {
	Src string `json:"src"` // What to copy (source).
	Dst string `json:"dst"` // Where to copy (destination).
}

// Themer manages themes.
type Themer struct {
	prefix string // Prefix directory where themes are stored.
}

// New returns a new instance of [Themer].
func New() *Themer {
	return NewWithPrefix(filepath.Join(xdg.ConfigHome, "theme"))
}

// NewWithPrefix is like [New] but allows to specify prefix directory.
func NewWithPrefix(prefix string) *Themer {
	return &Themer{prefix: prefix}
}

// Set sets a theme by name. It copies files from `copy.json` and runs `run`
// located in the theme directory.
// Returned slice is a slice of errors occurred during copying and execution.
func (t *Themer) Set(name string) (errors []error) {
	c, err := t.Copy(name)
	if err != nil {
		errors = append(errors, err)
	}

	var wg sync.WaitGroup

	for _, e := range c {
		wg.Add(1)
		go func() {
			err := cp.Copy(t.ExpandPath(e.Src, name), t.ExpandPath(e.Dst, name))
			if err != nil {
				errors = append(errors, err)
			}
			wg.Done()
		}()
	}

	err = t.Run(name)
	if err != nil {
		errors = append(errors, err)
	}

	wg.Wait()
	return errors
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

// Run runs `run` located in the theme directory.
func (t *Themer) Run(name string) error {
	path := filepath.Join(t.prefix, name, runfile)
	cmd := exec.Command(path)

	cmd.Dir = filepath.Join(t.prefix, name)

	return cmd.Run()
}

// Copy reads `copy.json` located in the theme directory and returns
// [Copy] entries.
func (t *Themer) Copy(name string) (c []Copy, err error) {
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
