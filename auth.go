package auth

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mordredp/auth/provider"
)

// authenticator manages sessions and authentication providers.
type authenticator struct {
	sessions         *store
	cookieName       string
	maxSessionLength time.Duration
	lastCleanup      time.Time
	providers        []provider.Provider
}

// New initializes a new authenticator.
func NewAuthenticator(maxSessionSeconds int, options ...Option) *authenticator {

	a := authenticator{
		sessions:         NewStore(),
		cookieName:       uuid.NewString(),
		maxSessionLength: time.Duration(maxSessionSeconds) * time.Second,
		lastCleanup:      time.Now(),
		providers:        make([]provider.Provider, 0),
	}

	for _, option := range options {
		err := option(&a)
		if err != nil {
			log.Printf("options: %s", err)
			continue
		}
	}

	go a.sessions.loop() // start a goroutine to manage the sessions store

	return &a
}
