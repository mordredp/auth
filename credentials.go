package auth

// Credentials holds the username and password used to authenticate a session.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
