package session

type SessionStore struct {
	Signing *SigningSessionStore

	// TODO: These only have 1 at a time
	// KeygenSessionStore *KeygenSessionStore
	// ReshareSessionStore *ReshareSessionStore
}

// NewSessionStore returns a new SessionStore.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		Signing: NewSigningSessionStore(),
	}
}
