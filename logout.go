package auth

import (
	"log"
	"net/http"
	"time"
)

// Logout removes a session.
func (a *authenticator) Logout(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie(a.cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			//w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusFound)

			return
		}

		//w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusFound)

		return
	}
	token := c.Value

	a.sessions.remove(token)

	log.Printf("session %q removed, %q logged out", token, r.Context().Value(UserKey).(User).ID)

	http.SetCookie(w, &http.Cookie{
		Name:     a.cookieName,
		Value:    "",
		Expires:  time.Now(),
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
