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
	providers        []provider.Provider
}

// New initializes a new authenticator.
func NewAuthenticator(maxSessionSeconds int, options ...Option) *authenticator {

	a := authenticator{
		sessions:         NewStore(),
		cookieName:       uuid.NewString(),
		maxSessionLength: time.Duration(maxSessionSeconds) * time.Second,
		providers:        make([]provider.Provider, 0),
	}

	for _, option := range options {
		err := option(&a)
		if err != nil {
			log.Printf("options: %s", err)
			continue
		}
	}

	go a.sessions.listen()                              // start a goroutine to listen for operations on the sessions store
	go a.sessions.startClearing(a.maxSessionLength / 2) // start a goroutine to clear expired sessions

	return &a
}
