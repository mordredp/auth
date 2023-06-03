package auth

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// An operation is a function that modifies a map of sessions.
type operation func(map[string]session)

// A store contains a channel that transmits operations.
type store struct {
	ops chan operation
}

// NewStore initializes the operations channel for the store.
func newStore() *store {
	return &store{
		ops: make(chan operation),
	}
}

// exists returns a session and 'true' if the token is
// associated to an authenticated session, otherwise it returns
// a zero value session and 'false'.
func (s *store) exists(token string) (session, bool) {
	result := make(chan session, 1)
	s.ops <- func(sessions map[string]session) {
		result <- sessions[token]
	}
	sess := <-result
	return sess, sess != session{}
}

// add adds a session to the store and returns a randomly generated
// token which identifies it.
func (s *store) add(sess session) string {
	result := make(chan string, 1)
	s.ops <- func(sessions map[string]session) {
		token := uuid.NewString()
		sessions[token] = sess
		result <- token
	}
	return <-result
}

// remove removes a session from the store, identified by a token.
// It sends an operation which deletes a session from the session map.
func (s *store) remove(token string) {
	s.ops <- func(sessions map[string]session) {
		delete(sessions, token)
	}
}

// clear removes expired sessions from the store. It sends an operation which
// deletes expired sessions from the session map and returns their count.
func (s *store) clear() int {
	result := make(chan int, 1)
	s.ops <- func(sessions map[string]session) {
		deletedCount := 0
		for token, session := range sessions {
			if session.expiry.Before(time.Now()) {
				delete(sessions, token)
				deletedCount++
			}
		}
		result <- deletedCount
	}
	return <-result
}

// startClearing starts a ticker with the specified period
// which will clear the store on every tick.
func (s *store) startClearing(period time.Duration) {
	log.Printf("store started clearing sessions every %s ...", period)
	ticker := time.NewTicker(period)

	for {
		tick := <-ticker.C
		clearedCount := s.clear()
		log.Printf("cleared %d sessions in %s", clearedCount, time.Since(tick))
	}
}

// listen initializes a map of sessions which associates a token to a session,
// starts listening on the channel for operations, and executes them.
func (s *store) listen() {
	log.Printf("store started listening for operations ...")
	sessions := make(map[string]session)
	for op := range s.ops {
		op(sessions)
	}
}
