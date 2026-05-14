package storage

import (
	"context"
	"errors"
	"sync"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (m *MemoryStore) Save(ctx context.Context, code, longUrl string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[code] = longUrl

	return nil
}

func (m *MemoryStore) Get(code string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	longUrl, exists := m.data[code]

	if !exists {
		return "", errors.New("Long url doesn't exist")
	}

	return longUrl, nil
}
