// Package errorlist provides a thread-safe way to manage a list of errors.
package errorlist

import "sync"

// Errorlist is a container that wraps `[]error`.
//
// Methods of Errorlist are safe to run in multiple goroutines.
type Errorlist struct {
	errors []error
	mu     sync.Mutex
}

// New returns a new instance of [Errorlist].
func New() *Errorlist {
	return &Errorlist{}
}

// Append appends error to the underlying errors slice.
func (l *Errorlist) Append(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.errors = append(l.errors, err)
}

// Get returns a copy of underlying errors slice. The returned value can be
// safely mutated.
func (l *Errorlist) Get() []error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Return a copy of the errors slice.
	copiedErrors := make([]error, len(l.errors))
	copy(copiedErrors, l.errors)
	return copiedErrors
}
