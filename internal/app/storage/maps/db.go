package maps

import (
	"errors"
	"sync"
)

type MapStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		data: make(map[string]string),
	}
}

func (m *MapStorage) Read(URL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[URL]
	if !ok {
		return "", errors.New("url not found")
	} else {
		return value, nil
	}
}

// создаем 2 записи: map[shortURL] = originURL, map[originURL] = shortURL
func (m *MapStorage) Create(shortURL, originURL string) string {
	val, err := m.Read(originURL)
	if err != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.data[shortURL] = originURL
		m.data[originURL] = shortURL
		return shortURL
	}
	return val
}

func (m *MapStorage) Delete(shortURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	originalURL, err := m.Read(shortURL)
	if err == nil {
		delete(m.data, originalURL)
		delete(m.data, shortURL)
		return nil
	}
	return errors.New("short url does not exist")
}
