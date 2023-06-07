package auth

import (
	"log"

	"github.com/mordredp/auth/provider"
	"github.com/mordredp/auth/provider/ldap"
	"github.com/mordredp/auth/provider/static"
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
		provider := static.Static(password)
		a.providers = append(a.providers, provider)
		log.Printf("configured %T", provider)

		return nil
	}
}
