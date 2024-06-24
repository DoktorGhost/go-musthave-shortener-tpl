package maps

import (
	"errors"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"sync"
)

type shortLong struct {
	shortURL string
	longURL  string
}

type MapStorage struct {
	random map[string]shortLong          //random:originalURL 	//shortURL:originalURL
	long   map[string]string             //long:short
	ids    map[string][]storage.UserURLs // id:[{short, long}, {short, long},...]
	mu     sync.RWMutex
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		random: make(map[string]shortLong),
		long:   make(map[string]string),
		ids:    make(map[string][]storage.UserURLs),
	}
}

// возвращаем по идентификатору сокращенной сылки "ыавыва" оригинальную ссылку  сокращенную ссылку
func (m *MapStorage) Read(short string) (string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.random[short]
	if !ok {
		return "", "", errors.New("url not found")
	} else {
		return value.longURL, value.shortURL, nil
	}
}

// получаем из оригинального УРЛа сокращенный идентификатор
func (m *MapStorage) ReadOriginal(originalURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.long[originalURL]
	if !ok {
		return "", errors.New("url not found")
	} else {
		return value, nil
	}
}

func (m *MapStorage) ReadUrlsUserId(userID string) ([]storage.UserURLs, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.ids[userID]
	if !ok {
		return nil, errors.New("not url with UserID = " + userID)
	} else {
		return value, nil
	}
}

func (m *MapStorage) Create(short, shortURL, originURL, userID string) error {
	if len(m.long[originURL]) < 1 {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.random[short] = shortLong{shortURL: shortURL, longURL: originURL}
		m.long[originURL] = short
		m.ids[userID] = append(m.ids[userID], storage.UserURLs{ShortURL: shortURL, LongURL: originURL})
		return nil
	} else {
		return usecase.ErrShortURLAlreadyExists
	}
}
