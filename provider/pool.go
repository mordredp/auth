package provider

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// A pool is a set of Providers
type Pool []Provider

type result struct {
	provider *Provider
	err      error
}

// Authenticate calls starts an authentication goroutine on all Providers in the
// pool in a concurrent fashion. It waits for all Providers to complete and
// returns the first Provider which authenticated successfully or an error.
func (p Pool) Authenticate(creds Credentials) (*Provider, error) {

	results := make(chan result, len(p))

	var wg sync.WaitGroup
	start := time.Now()

	for _, prov := range p {
		wg.Add(1)
		go func(p Provider) {
			results <- result{
				err:      p.Authenticate(creds),
				provider: &p,
			}
			wg.Done()
		}(prov)
	}

	wg.Wait()
	close(results)

	elapsed := time.Since(start)
	jitter, err := rand.Int(rand.Reader, big.NewInt(elapsed.Nanoseconds()/8))
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(jitter.Int64()))

	for result := range results {
		if result.err == nil {
			return result.provider, result.err
		}
	}

	return nil, fmt.Errorf("no providers authenticated %q", creds.Username)
}
