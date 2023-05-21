WIP README

# auth
`auth` is a library written in Go to provide authentication for the standard HTTP server.
It strives to be idiomatic, depend only on stdlib, and be simple to use.

## Install
`go get -u github.com/mordredp/auth`

## Features
* **100% compatible with net/http** - use any http or middleware pkg in the ecosystem that is also compatible with `net/http`
<!---
* **No external dependencies** - plain ol' Go stdlib + `net/http`
-->
* **context** - built on new `context` package to pass information to the handlers
* **go.mod support** - 

## Examples
```go
package main

import (
	"net/http"

	"github.com/mordredp/auth"
)

func main() {
	authenticator := auth.New(120, auth.Static("test"))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

    //TODO: add example

	http.ListenAndServe(":3000", mux)
}

```

## Handlers
### Core handlers
------------------------------------------------------
| Handler                | description               |
| :--------------------- | :------------------------ |
| [auth.Login]           | creates a session         |
| [auth.Logout]          | deletes a session         |
| [auth.Refresh]         | refreshes a session token |
| [TODO]                 | TODO                      |
------------------------------------------------------

## Providers
The library currently implements a `Default` provider (which never authenticates), a `Static` provider to be used for testing and an `LDAP` provider.


