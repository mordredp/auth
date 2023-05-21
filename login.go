package auth

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Login authenticates the session assigned to a user.
// It tries to authenticate the session on all providers configured,
// and returns as soon as the first one succeeds.
func (a *authenticator) Login(w http.ResponseWriter, r *http.Request) {

	c := Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var err error = fmt.Errorf("no providers authenticated user %q", c.Username)
	var goodProvider string

	for _, provider := range a.providers {

		if err := provider.Authenticate(c.Username, c.Password); err != nil {
			log.Printf("provider: %s", err.Error())
			continue
		}

		goodProvider = reflect.TypeOf(provider).String()
		err = nil
		break
	}

	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(a.maxSessionLength)

	a.sessions[sessionToken] = session{
		ID:     c.Username,
		expiry: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     a.cookieName,
		Value:    sessionToken,
		Expires:  expiresAt,
		SameSite: http.SameSiteStrictMode,
	})

	log.Printf("user %q logged in with %q", c.Username, goodProvider)

	http.Redirect(w, r, "/", http.StatusFound)
}
