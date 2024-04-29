package maps

import (
	"errors"
	"fmt"
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

func (m *MapStorage) Read(shortURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fmt.Println(m.data[shortURL])
	value, ok := m.data[shortURL]
	if !ok {
		return "", errors.New("short url already exists")
	} else {
		return value, nil
	}
}

// создаем 2 записи: map[shortURL] = originURL, map[originURL] = shortURL
func (m *MapStorage) Create(shortURL, originURL string) error {
	_, err := m.Read(originURL)
	if err != nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.data[shortURL] = originURL
		m.data[originURL] = shortURL
		return nil
	} else {
		return errors.New("short url already exists")
	}
}

func (m *MapStorage) Delete(shortURL string) error {
	originalURL, err := m.Read(shortURL)
	if err == nil {
		delete(m.data, originalURL)
		delete(m.data, shortURL)
		return nil
	}
	return errors.New("short url does not exist")
}
