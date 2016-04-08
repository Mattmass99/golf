package Golf

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const sessionIDLength = 64

// MemorySessionManager manages a map of sessions.
type MemorySessionManager struct {
	sessions map[string]*MemorySession
}

// NewMemorySessionManager creates a new session manager.
func NewMemorySessionManager() *MemorySessionManager {
	mgr := new(MemorySessionManager)
	mgr.sessions = make(map[string]*MemorySession)
	return mgr
}

func (mgr *MemorySessionManager) sessionID() (string, error) {
	b := make([]byte, sessionIDLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

// Session gets the session instance by indicating a session id.
func (mgr *MemorySessionManager) Session(sid string) (Session, error) {
	if s, ok := mgr.sessions[sid]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Can not retrieve session with id %s.", sid)
}

// NewSession creates and returns a new session.
func (mgr *MemorySessionManager) NewSession() (Session, error) {
	sid, err := mgr.sessionID()
	if err != nil {
		return nil, err
	}
	s := MemorySession{sid: sid, data: make(map[string]interface{})}
	mgr.sessions[sid] = &s
	return &s, nil
}

// Session is an interface for session instance, a session instance contains
// data needed.
type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	SessionID() string
}

// MemorySession is an memory based implementaion of Session.
type MemorySession struct {
	sid  string
	data map[string]interface{}
}

// Set method sets the key value pair in the session.
func (s *MemorySession) Set(key string, value interface{}) error {
	s.data[key] = value
	return nil
}

// Get method gets the value by given a key in the session.
func (s *MemorySession) Get(key string) (interface{}, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("key %q in session (id %d) not found", key, s.sid)
}

// Delete method deletes the value by given a key in the session.
func (s *MemorySession) Delete(key string) error {
	delete(s.data, key)
	return nil
}

// SessionID returns the current ID of the session.
func (s *MemorySession) SessionID() string {
	return s.sid
}