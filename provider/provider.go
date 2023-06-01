package provider

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
