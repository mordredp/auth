package main

import (
	"net/http"
	"text/template"

	"github.com/mordredp/auth"
)

func main() {
	authenticator := auth.NewAuthenticator(30, auth.Static("test"))

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
	mux.HandleFunc("/login", authenticator.Login)
	mux.HandleFunc("/logout", authenticator.Logout)

	http.ListenAndServe(":8080", mux)
}
