package auth

import (
	"time"
)

type key int

const (
	// UserKey is the key to the value of a User in a context.
	UserKey key = iota
)

// User holds a users account information and its authentication status.
type User struct {
	ID            string
	Authenticated bool
}

// Credentials holds the username and password used to authenticate a session.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// a session contains an identifier (usually the username of the user
// it is assigned to) and an expiration time.
type session struct {
	ID     string
	expiry time.Time
}

// isExpired returns true if the session has expired,
// false otherwise.
func (s *session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// userOrDefault returns an authenticated User
// from the Authenticator session store, or an empty one.
func (a *authenticator) userOrDefault(id string) User {
	session, ok := a.sessions[id]
	if !ok {
		return User{}
	}
	return User{ID: session.ID, Authenticated: true}
}
