package auth

import (
	"log"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/mordredp/auth/provider"
)

// authenticator manages sessions and authentication providers.
type authenticator struct {
	sessions         map[string]session
	cookieName       string
	maxSessionLength time.Duration
	lastCleanup      time.Time
	tpl              *template.Template
	providers        []provider.Provider
}

// An Option modifies an authenticator or returns an error.
type Option func(*authenticator) error

// New initializes a new authenticator.
func New(sessionSeconds int, options ...Option) *authenticator {

	a := authenticator{
		sessions:         make(map[string]session),
		cookieName:       uuid.NewString(),
		maxSessionLength: time.Duration(sessionSeconds) * time.Second,
		lastCleanup:      time.Now(),
		tpl:              template.Must(template.ParseGlob("auth/templates/*.gohtml")),
		providers:        make([]provider.Provider, 0),
	}

	for _, option := range options {
		err := option(&a)
		if err != nil {
			log.Printf("options: %s", err)
			continue
		}
	}

	return &a
}
