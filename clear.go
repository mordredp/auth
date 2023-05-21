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

			deletedCount := 0
			for id, session := range a.sessions {
				if session.isExpired() {
					delete(a.sessions, id)
					deletedCount++
				}
			}
			a.lastCleanup = time.Now()

			log.Printf("removed %d expired sessions", deletedCount)
		}

		next.ServeHTTP(w, r)
	})
}
