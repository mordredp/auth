package auth

import (
	"time"
)

// A session contains an identifier (usually the username of the user
// it is assigned to) and an expiration time.
type session struct {
	ID     string
	expiry time.Time
}
