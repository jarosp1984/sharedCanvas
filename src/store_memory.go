package main

import (
	"context"
	"sync"
)

type memorySessionState struct {
	lines  []Line
	nextID int
}

type MemoryLineStore struct {
	mu       sync.Mutex
	sessions map[string]*memorySessionState
}

func NewMemoryLineStore() *MemoryLineStore {
	return &MemoryLineStore{
		sessions: make(map[string]*memorySessionState),
	}
}

func (store *MemoryLineStore) sessionState(sessionID string) *memorySessionState {
	state, exists := store.sessions[sessionID]
	if exists {
		return state
	}

	state = &memorySessionState{
		nextID: 1,
	}
	store.sessions[sessionID] = state

	return state
}

func (store *MemoryLineStore) nextLineIDLocked(sessionID string) int {
	state := store.sessionState(sessionID)
	id := state.nextID
	state.nextID++

	return id
}

func (store *MemoryLineStore) NextLineID(_ context.Context, sessionID string) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	return store.nextLineIDLocked(sessionID), nil
}

func (store *MemoryLineStore) SaveLine(_ context.Context, sessionID string, line *Line) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	line.ID = store.nextLineIDLocked(sessionID)
	state := store.sessionState(sessionID)

	state.lines = append(state.lines, *line)

	return nil
}

func (store *MemoryLineStore) ListLines(_ context.Context, sessionID string) ([]Line, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	state := store.sessionState(sessionID)
	lines := make([]Line, len(state.lines))
	copy(lines, state.lines)

	return lines, nil
}

func (store *MemoryLineStore) ListLinesSince(_ context.Context, sessionID string, sinceID int) ([]Line, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	state := store.sessionState(sessionID)
	lines := make([]Line, 0, len(state.lines))
	for _, line := range state.lines {
		if line.ID > sinceID {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func (store *MemoryLineStore) ClearLines(_ context.Context, sessionID string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	state, exists := store.sessions[sessionID]
	if !exists {
		return nil
	}

	state.lines = nil

	return nil
}
