package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/mordredp/auth/provider"
)

// Login authenticates the session assigned to a user.
// It tries to authenticate the session on all configured providers
// in order of registration and returns as soon as the first one succeeds.
func (a *authenticator) Login(w http.ResponseWriter, r *http.Request) {

	creds := provider.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	provider, err := a.providers.Authenticate(creds)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	log.Printf("%T authenticated %q", *provider, creds.Username)

	expiresAt := time.Now().Add(a.maxSessionLength)

	token := a.sessions.add(session{
		ID:     creds.Username,
		expiry: expiresAt,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     a.cookieName,
		Value:    token,
		Expires:  expiresAt,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
