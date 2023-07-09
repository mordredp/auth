package ldap

import (
	"crypto/tls"

	"github.com/go-ldap/ldap"
	"github.com/mordredp/auth/provider"
)

// An Option modifies a directory or returns an error.
type Option func(d *directory) error

// Bind verifies both the connection and bind status
// to a directory with the credentials provided to it.
func Bind(creds provider.Credentials) func(d *directory) error {
	return func(d *directory) error {
		conn, err := ldap.DialURL(d.bindAddr.String())
		if err != nil {
			return err
		}

		defer conn.Close()

		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}

		d.bindUser = creds.Username
		d.bindPass = creds.Password

		return conn.Bind(d.bindUser, d.bindPass)
	}
}

// IdKey sets the name of the query parameter to use for identification.
func IdKey(key string) func(d *directory) error {
	return func(d *directory) error {
		d.idKey = key

		return nil
	}
}

// QueryParams sets the query parameters to use for identification
func QueryParams(params map[string]string) func(d *directory) error {
	return func(d *directory) error {
		for key, value := range params {
			d.queryParams[key] = value
		}

		return nil
	}
}
