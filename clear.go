package auth

import (
	"log"
	"net/http"
	"time"
)

// Clear removes expired sessions.
func (a *authenticator) Clear(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if time.Now().After(a.lastCleanup.Add(a.maxSessionLength / 2)) {

			clearedCount := a.sessions.clear()
			a.lastCleanup = time.Now()

			log.Printf("removed %d expired sessions", clearedCount)
		}

		next.ServeHTTP(w, r)
	})
}
