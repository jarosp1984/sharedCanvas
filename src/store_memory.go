package main

import (
	"context"
	"sync"
)

type MemoryLineStore struct {
	mu     sync.Mutex
	lines  []Line
	nextID int
}

func NewMemoryLineStore() *MemoryLineStore {
	return &MemoryLineStore{
		nextID: 1,
	}
}

func (store *MemoryLineStore) NextLineID(_ context.Context) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	id := store.nextID
	store.nextID++

	return id, nil
}

func (store *MemoryLineStore) SaveLine(ctx context.Context, line *Line) error {
	id, err := store.NextLineID(ctx)
	if err != nil {
		return err
	}

	line.ID = id

	store.mu.Lock()
	defer store.mu.Unlock()

	store.lines = append(store.lines, *line)

	return nil
}

func (store *MemoryLineStore) ListLines(_ context.Context) ([]Line, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	lines := make([]Line, len(store.lines))
	copy(lines, store.lines)

	return lines, nil
}

func (store *MemoryLineStore) ListLinesSince(_ context.Context, sinceID int) ([]Line, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	lines := make([]Line, 0, len(store.lines))
	for _, line := range store.lines {
		if line.ID > sinceID {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
