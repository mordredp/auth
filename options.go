package auth

import (
	"log"

	"github.com/mordredp/auth/provider"
	"github.com/mordredp/auth/provider/ldap"
)

// An Option modifies an authenticator or returns an error.
type Option func(*authenticator) error

func CookieName(key string) Option {
	return func(a *authenticator) error {
		a.cookieName = key

		return nil
	}
}

// LDAP adds an LDAP provider to the Authenticator.
func LDAP(addr string, baseDN string, username string, password string, options ...ldap.Option) Option {
	return func(a *authenticator) error {

		ldapOptions := append([]ldap.Option{}, options...)
		ldapOptions = append(ldapOptions, ldap.Bind(username, password))

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
		static := provider.Static(password)
		a.providers = append(a.providers, static)
		log.Printf("configured %T", static)

		return nil
	}
}
