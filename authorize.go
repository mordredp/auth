package auth

import "net/http"

// Authorize retrieves a User from a context, and it passes control to the
// next handler in the chain if the User is Authenticated.
func (a *authenticator) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(User)

		if !user.Authenticated {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
