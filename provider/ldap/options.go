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

// Fields sets some parameters for the LDAP filter:
// classValue sets the value for the parameter "objectClass"
// and idKey sets the name of the field to use for identification.
func Fields(classValue string, idKey string) func(d *directory) error {
	return func(d *directory) error {
		d.classValue = classValue
		d.idKey = idKey

		return nil
	}
}
