package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Identify retrieves a session or creates a new one.
func (a *authenticator) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie(a.cookieName)

		var sessionToken string

		switch err {
		case nil:
			sessionToken = c.Value

		case http.ErrNoCookie:
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(a.maxSessionLength)

			http.SetCookie(w, &http.Cookie{
				Name:     a.cookieName,
				Value:    sessionToken,
				Expires:  expiresAt,
				SameSite: http.SameSiteStrictMode,
			})

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := a.userOrDefault(sessionToken)

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
