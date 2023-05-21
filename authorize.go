package auth

import "net/http"

func (a *authenticator) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(User)

		if !user.Authenticated {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	})
}
