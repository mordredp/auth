package auth

import (
	"log"

	"github.com/mordredp/auth/provider"
	"github.com/mordredp/auth/provider/ldap"
)

// LDAP adds an LDAP provider to the Authenticator.
func LDAP(addr string, baseDN string, username string, password string, options ...ldap.Option) func(a *authenticator) error {
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
		log.Printf("configured LDAP provider on %q with base DN %q", addr, baseDN)

		return nil
	}
}

// Static adds a Static provider to the Authenticator.
func Static(password string) func(a *authenticator) error {
	return func(a *authenticator) error {
		a.providers = append(a.providers, provider.Static(password))
		log.Printf("configured Static provider ")

		return nil
	}
}
