package main

import (
	"net/http"
	"text/template"

	"github.com/mordredp/auth"
)

func main() {
	authenticator := auth.NewAuthenticator(auth.Static("test"))

	mux := http.NewServeMux()

	home := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var user auth.User
		if v, ok := r.Context().Value(auth.UserKey).(auth.User); ok {
			user = v
		}

		if !user.Authenticated {
			tpl := template.Must(template.ParseGlob("*.gohtml"))
			tpl.ExecuteTemplate(w, "index.gohtml", user)
		} else {
			w.Write([]byte("welcome " + user.ID))
		}
	})

	mux.Handle("/", authenticator.Identify(home))
	mux.Handle("/login", authenticator.Identify(http.HandlerFunc(authenticator.Login)))
	mux.Handle("/logout", authenticator.Identify(http.HandlerFunc(authenticator.Logout)))

	http.ListenAndServe(":8080", mux)
}
