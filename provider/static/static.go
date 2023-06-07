package static

import (
	"github.com/mordredp/auth/provider"
	"github.com/pkg/errors"
)

// Static is an implementation of a Provider.
type Static string

// Authenticate returns an error if the provided password does not match
// the one Static has been set to
func (s Static) Authenticate(creds provider.Credentials) error {
	if creds.Password != string(s) {
		return errors.New("invalid password")
	}

	return nil
}
