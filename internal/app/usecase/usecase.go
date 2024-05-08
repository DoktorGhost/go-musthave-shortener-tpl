package usecase

import (
	"errors"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/shortener"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage"
	"net/url"
)

type ShortURLUseCase struct {
	storage storage.Repository
}

func NewShortURLUseCase(storage storage.Repository) *ShortURLUseCase {
	return &ShortURLUseCase{storage: storage}
}

func (uc *ShortURLUseCase) CreateShortURL(originalURL string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", err
	}
	
	for i := 0; i < 10; i++ {
		shortURL := shortener.RandomString(8)
		_, err := uc.storage.Read(shortURL)
		if err != nil {
			return uc.storage.Create(shortURL, originalURL), nil
		}
	}
	return "", errors.New("short url already exists")
}

func (uc *ShortURLUseCase) GetOriginalURL(shortURL string) (string, error) {
	originalURL, err := uc.storage.Read(shortURL)
	if err != nil {
		return "", err
	} else {
		return originalURL, nil
	}
}

func (uc *ShortURLUseCase) DeleteURL(shortURL string) error {
	return uc.storage.Delete(shortURL)
}
