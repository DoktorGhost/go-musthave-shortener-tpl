package usecase

import (
	"errors"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/shortener"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage"
)

type ShortUrlUseCase struct {
	storage storage.Repository
}

func NewShortUrlUseCase(storage storage.Repository) *ShortUrlUseCase {
	return &ShortUrlUseCase{storage: storage}
}

func (uc *ShortUrlUseCase) CreateShortUrl(originalURL string) (string, error) {
	for i := 0; i < 10; i++ {
		shortURL := shortener.RandomString(8)
		_, err := uc.storage.Read(shortURL)
		if err != nil {
			return uc.storage.Create(shortURL, originalURL), nil
		}
	}
	return "", errors.New("short url already exists")
}

func (uc *ShortUrlUseCase) GetOriginalUrl(shortURL string) (string, error) {
	originalURL, err := uc.storage.Read(shortURL)
	if err != nil {
		return "", err
	} else {
		return originalURL, nil
	}
}

func (uc *ShortUrlUseCase) DeleteUrl(shortURL string) error {
	return uc.storage.Delete(shortURL)
}
