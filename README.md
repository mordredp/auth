WIP README

# auth
`auth` is a library written in Go to provide authentication for the standard HTTP server.
It strives to be idiomatic, depend mostly on stdlib, and be simple to use.

## Install
`go get -u github.com/mordredp/auth`

## Features
* **100% compatible with net/http** - use any http or middleware pkg in the ecosystem that is also compatible with `net/http`
<!---
* **No external dependencies** - plain ol' Go stdlib + `net/http`
-->
* **context** - built on new `context` package to pass information to the handlers
* **go.mod support** - go.mod which lists all dependencies included
* **extensible** - `Provider` interface allows to implement custom authentication providers (`Default`, `Static` and `LDAP` already included)

## Examples
```go
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
```

## Handlers
### Core handlers
------------------------------------------------------
| Handler                | description               |
| :--------------------- | :------------------------ |
| [auth.Login]           | creates a session         |
| [auth.Logout]          | deletes a session         |
| [auth.Authorize]       | authorizes a session      |
------------------------------------------------------

## Providers
The library currently implements a `Default` provider (which never authenticates), a `Static` provider to be used for testing and an `LDAP` provider.


