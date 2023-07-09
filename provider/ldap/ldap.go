package ldap

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-ldap/ldap"
	"github.com/mordredp/auth/provider"
	"github.com/pkg/errors"
)

// A directory represents an LDAP search domain.
type directory struct {
	bindAddr    url.URL
	bindUser    string
	bindPass    string
	baseDN      string
	idKey       string
	queryParams map[string]string
}

// NewDirectory initializes an ldap client. The initialization fails if any
// option returns an error.
func NewDirectory(addr string, baseDN string, options ...Option) (*directory, error) {

	url, err := url.Parse(addr)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%q is an invalid address", addr))
	}

	d := directory{
		bindAddr:    *url,
		baseDN:      baseDN,
		idKey:       "id",
		queryParams: make(map[string]string),
	}

	for _, option := range options {
		err := option(&d)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

// buildQuery returns a query string to be used for authentication. It builds
// the query using the idKey and queryParams set on the directory.
func (d *directory) buildQuery(creds provider.Credentials) string {

	d.queryParams[d.idKey] = ldap.EscapeFilter(creds.Username)

	paramFmt := "(%s=%s)"
	var queryFmtBuilder strings.Builder

	var andParams []any

	for key, value := range d.queryParams {
		andParams = append(andParams, fmt.Sprintf(paramFmt, key, value))
		queryFmtBuilder.WriteString("%s")
	}

	return fmt.Sprintf("(&%s)", fmt.Sprintf(queryFmtBuilder.String(), andParams...))
}

// Authenticate returns an error if the username is not found within
// the directory or the username does not bind to it with the provided password.
func (d *directory) Authenticate(creds provider.Credentials) error {

	conn, err := ldap.DialURL(d.bindAddr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}

	err = conn.Bind(d.bindUser, d.bindPass)
	if err != nil {
		return err
	}

	query := d.buildQuery(creds)

	searchRequest := ldap.NewSearchRequest(
		d.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		query,
		[]string{"dn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return errors.Errorf("found %d entries for %q", len(sr.Entries), query)
	}

	userdn := sr.Entries[0].DN

	err = conn.Bind(userdn, creds.Password)
	if err != nil {
		return err
	}

	return nil
}
