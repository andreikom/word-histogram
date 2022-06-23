package storage

import "sync"

type WordStore interface {
	Add(map[string]int)
	Get() map[string]int
}

type InMemory struct {
	store map[string]int
	lock  sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		store: make(map[string]int, 0),
	}
}

func (m *InMemory) Add(wordCount map[string]int) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for word, newCount := range wordCount {
		count, ok := m.store[word]
		if ok {
			m.store[word] = count + newCount

			continue
		}

		m.store[word] = newCount
	}
}

func (m *InMemory) Get() map[string]int {
	storeCopy := make(map[string]int)
	m.lock.RLock()
	defer m.lock.RUnlock()

	for key, value := range m.store {
		storeCopy[key] = value
	}

	return storeCopy
}
