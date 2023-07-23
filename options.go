package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mordredp/auth/provider"
	"github.com/mordredp/auth/provider/ldap"
	"github.com/mordredp/auth/provider/static"
)

// An Option modifies an authenticator or returns an error.
type Option func(*authenticator) error

// MaxSessionLength sets the maximum session length. If the number of seconds is
// not valid the default duration will be used.
func MaxSessionLength(seconds int) Option {
	return func(a *authenticator) error {

		if seconds < 1 {
			return fmt.Errorf("invalid session duration specified, using %s", a.maxSessionLength)
		}

		a.maxSessionLength = time.Duration(seconds) * time.Second

		return nil
	}
}

// CookieName sets the name of the cookie to use for session identification. If
// the cookie name is not valid a default name will be used.
func CookieName(key string) Option {
	return func(a *authenticator) error {

		if key == "" {
			return fmt.Errorf("invalid cookie name specified, using %s", a.cookieName)
		}

		a.cookieName = key

		return nil
	}
}

// LDAP adds an LDAP provider to the Authenticator.
func LDAP(addr string, baseDN string, creds provider.Credentials, options ...ldap.Option) Option {
	return func(a *authenticator) error {

		ldapOptions := append([]ldap.Option{}, options...)
		ldapOptions = append(ldapOptions, ldap.Bind(provider.Credentials{
			Username: creds.Username,
			Password: creds.Password,
		}))

		ldap, err := ldap.NewDirectory(
			addr,
			baseDN,
			ldapOptions...,
		)

		if err != nil {
			return err
		}

		a.providers = append(a.providers, ldap)
		log.Printf("configured %T on %q with base DN %q", ldap, addr, baseDN)

		return nil
	}
}

// Static adds a Static provider to the Authenticator.
func Static(password string) Option {
	return func(a *authenticator) error {

		if password == "" {
			return errors.New("static password not specified")
		}

		provider := static.Static(password)
		a.providers = append(a.providers, provider)

		log.Printf("configured %T", provider)

		return nil
	}
}
