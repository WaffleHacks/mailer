package daemon

import (
	"github.com/gobwas/glob"

	"github.com/WaffleHacks/mailer/providers"
)

type Matcher struct {
	id       string
	match    glob.Glob
	provider providers.Provider
	queue    chan Message
	workers  int
}

// NewMatcher creates a new matcher for the provider with its associated queue
func NewMatcher(id string, numWorkers int, provider providers.Provider, pattern string) (*Matcher, error) {
	// Ensure the glob compiles
	g, err := glob.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Matcher{
		id:       id,
		match:    g,
		provider: provider,
		queue:    make(chan Message),
		workers:  numWorkers,
	}, nil
}

// Enqueue adds the message to the provider's queue
func (m *Matcher) Enqueue(msg Message) {
	if m.match.Match(msg.From) {
		m.queue <- msg
	}
}
