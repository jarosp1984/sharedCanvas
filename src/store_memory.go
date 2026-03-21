package main

import (
	"context"
	"sync"
)

type MemoryLineStore struct {
	mu    sync.Mutex
	lines []Line
}

func NewMemoryLineStore() *MemoryLineStore {
	return &MemoryLineStore{}
}

func (store *MemoryLineStore) SaveLine(_ context.Context, line Line) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.lines = append(store.lines, line)

	return nil
}

func (store *MemoryLineStore) ListLines(_ context.Context) ([]Line, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	lines := make([]Line, len(store.lines))
	copy(lines, store.lines)

	return lines, nil
}
