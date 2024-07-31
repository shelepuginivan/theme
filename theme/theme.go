// Package theme provides theme management capabilities.
package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	quiet  bool   // Whether to supress warnings and subprocess output.
}

// Config is a configuration for [Themer].
type Config struct {
	Prefix string // Prefix directory where themes are stored.
	Quiet  bool   // Whether to supress warnings and subprocess output.
}

// New returns a new instance of [Themer].
func New() *Themer {
	return &Themer{
		prefix: filepath.Join(xdg.ConfigHome, "theme"),
	}
}

// NewWithConfig is like [New] but allows to specify configuration.
func NewWithConfig(cfg Config) *Themer {
	return &Themer{
		prefix: cfg.Prefix,
		quiet:  cfg.Quiet,
	}
}

// Set sets a theme by name. It copies files from `copy.json` and runs `run`
// located in the theme directory.
func (t *Themer) Set(name string) {
	c, err := t.ReadCopyfile(name)
	if err != nil {
		t.Warn(fmt.Sprintf("%s not found", copyfile))
	}

	for _, e := range c {
		err := cp.Copy(t.ExpandPath(e.Src, name), t.ExpandPath(e.Dst, name))
		if err != nil {
			t.Warn(fmt.Sprintf("cannot copy %s to %s: %s", e.Src, e.Dst, err))
		}
	}

	err = t.ExecRunfile(name)
	if err != nil {
		t.Warn(fmt.Sprintf("cannot run %s: %s", runfile, err))
	}
}

// List prints available themes to stdout.
func (t *Themer) List() {
	entries, err := os.ReadDir(t.prefix)
	if err != nil {
		t.Warn(fmt.Sprintf("cannot list themes: %s", err))
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			fmt.Println(e.Name())
		}
	}
}

// ExecRunfile runs `run` located in the theme directory.
func (t *Themer) ExecRunfile(name string) error {
	path := filepath.Join(t.prefix, name, runfile)
	cmd := exec.Command(path)

	cmd.Dir = filepath.Join(t.prefix, name)

	if !t.quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// ReadCopyfile reads `copy.json` located in the theme directory and returns
// [Copy] entries.
func (t *Themer) ReadCopyfile(name string) (c []Copy, err error) {
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

// Warn prints warning message to stderr.
func (t *Themer) Warn(msg string) {
	if t.quiet {
		return
	}
	fmt.Fprintln(os.Stderr, msg)
}
