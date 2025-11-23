package store

import (
	"sync"

	"github.com/gitglubber/email-client/internal/models"
)

// Simple in-memory store for sessions
// In production, use Redis or a proper session store
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*models.User
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*models.User),
	}
}

func (s *SessionStore) Set(sessionID string, user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = user
}

func (s *SessionStore) Get(sessionID string) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.sessions[sessionID]
	return user, exists
}

func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}
