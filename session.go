package auth

import (
	"time"

	"github.com/google/uuid"
)

// A session contains an identifier (usually the username of the user
// it is assigned to) and an expiration time.
type session struct {
	ID     string
	expiry time.Time
}

// A store contains a channel that receives functions which operate on a map.
type store struct {
	ops chan func(map[string]session)
}

// NewStore initializes the sessions store channel that will receive
// the operations to do on a map.
func NewStore() *store {
	return &store{
		ops: make(chan func(map[string]session)),
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
func (s *store) remove(token string) {
	s.ops <- func(sessions map[string]session) {
		delete(sessions, token)
	}
}

// clear removes expired sessions from the store.
func (s *store) clear() int {
	result := make(chan int, 1)
	s.ops <- func(sessions map[string]session) {
		deleted := 0
		for token, session := range sessions {
			if session.expiry.Before(time.Now()) {
				delete(sessions, token)
				deleted++
			}
		}
		result <- deleted
	}
	return <-result
}

// loop initializes the sessions map which associates a token to a session,
// and starts listening on the store channel for operations.
func (s *store) loop() {
	sessions := make(map[string]session)
	for op := range s.ops {
		op(sessions)
	}
}
