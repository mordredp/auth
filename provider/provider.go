package provider

import "github.com/pkg/errors"

// Credentials holds the username and password used to Authenticate on a provider.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// A Provider can Authenticate Credentials.
type Provider interface {
	// Authenticate returns an error if the username and password are not valid.
	Authenticate(Credentials) error
}

// Default is the default implementation of a Provider.
type Default struct{}

// Authenticate always returns an error signaling that
// the default provider is being used.
func (d *Default) Authenticate(creds Credentials) error {
	return errors.New("will never authenticate")
}
