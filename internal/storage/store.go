package storage

import (
	"sync"
)

type Store struct {
	Mu      sync.RWMutex
	Data    map[string]string
	Counter int64
}

func NewStore() *Store {
	return &Store{
		Data: make(map[string]string),
	}
}
