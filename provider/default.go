package provider

import (
	"errors"
)

// Default is a default implementation of a Provider.
type Default struct{}

// Authenticate always returns an error signaling that
// the default provider is being used.
func (d *Default) Authenticate(username string, password string) error {
	return errors.New("will never authenticate")
}
