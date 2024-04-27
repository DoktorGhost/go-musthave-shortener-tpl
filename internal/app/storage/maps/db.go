package maps

import "sync"

type MapStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		data: make(map[string]string),
	}
}

func (m *MapStorage) Create(shortURL, originURL string) error{
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.data[shortURL] {
		if !ok {
			return errors.New("short url already exists")
		}
	}
	m.data[shortURL] = originURL
	return nil
}