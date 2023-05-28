package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Login authenticates the session assigned to a user.
// It tries to authenticate the session on all providers configured,
// and returns as soon as the first one succeeds.
func (a *authenticator) Login(w http.ResponseWriter, r *http.Request) {

	c := Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var err error = fmt.Errorf("no providers authenticated %q", c.Username)

	for _, provider := range a.providers {

		if err := provider.Authenticate(c.Username, c.Password); err != nil {
			log.Printf("provider: %s", err.Error())
			continue
		}

		log.Printf("%T authenticated %q", provider, c.Username)
		err = nil
		break
	}

	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	expiresAt := time.Now().Add(a.maxSessionLength)

	token := a.sessions.add(session{
		ID:     c.Username,
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
