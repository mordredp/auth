package provider

import (
	"github.com/pkg/errors"
)

// Default is a default implementation of a Provider.
type Default struct{}

// Authenticate always returns an error signaling that
// the default provider is being used.
func (d *Default) Authenticate(creds Credentials) error {
	return errors.New("will never authenticate")
}
