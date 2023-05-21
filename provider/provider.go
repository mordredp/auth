package provider

// A Provider can Authenticate a pair of username and password.
type Provider interface {
	// Authenticate returns an error if the username and password are not valid.
	Authenticate(username string, password string) error
}
