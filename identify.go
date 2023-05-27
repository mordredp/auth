package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type key int

const (
	// UserKey is the key to the value of a User in a context.
	UserKey key = iota
)

// User holds a users account information and its authentication status.
type User struct {
	ID            string
	Authenticated bool
}

// Identify retrieves a session from the store or creates a new anonymous one.
// It stores a User in the context.
func (a *authenticator) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie(a.cookieName)

		var token string

		switch err {
		case nil:
			token = c.Value

		case http.ErrNoCookie:
			token := uuid.NewString()
			expiresAt := time.Now().Add(a.maxSessionLength)

			http.SetCookie(w, &http.Cookie{
				Name:     a.cookieName,
				Value:    token,
				Expires:  expiresAt,
				SameSite: http.SameSiteStrictMode,
			})

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		usr := User{}

		sess, ok := a.sessions.exists(token)
		if ok {
			usr.ID = sess.ID
			usr.Authenticated = true
		}

		ctx := context.WithValue(r.Context(), UserKey, usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
