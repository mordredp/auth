package provider

import (
	"errors"
)

// Static is an implementation of a Provider.
type Static string

// Authenticate returns an error if the provided password does not match
// the one Static has been set to
func (s Static) Authenticate(username string, password string) error {
	if password != string(s) {
		return errors.New("Static: invalid password")
	}

	return nil
}
