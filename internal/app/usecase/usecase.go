package usecase

import (
	"errors"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/osfile"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/shortener"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage"
	"log"
	"net/url"
	"strconv"
	"time"
)

type ShortURLUseCase struct {
	storage storage.Repository
}

func NewShortURLUseCase(storage storage.Repository) *ShortURLUseCase {
	return &ShortURLUseCase{storage: storage}
}

var ErrShortURLAlreadyExists = errors.New("short url already exists")

func (uc *ShortURLUseCase) CreateShortURL(originalURL string, conf *config.Config) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", err
	}

	//рандомная строка, будующая ссылка
	shortURL := shortener.RandomString(8)

	err = uc.storage.Create(shortURL, originalURL)
	//запись в файл
	if err == nil {
		if conf.FileStoragePath != "" {
			prod, err := osfile.NewProducer(conf.FileStoragePath)
			if err != nil {
				log.Printf("Ошибка создания Producer: %v\n", err)
				return shortURL, nil
			} else {
				currentTime := time.Now()
				intFromTime := currentTime.Unix()
				event := osfile.Event{
					UUID:        strconv.Itoa(int(intFromTime)),
					ShortURL:    shortURL,
					OriginalURL: originalURL,
				}
				err = prod.WriteEvent(&event)
				if err != nil {
					log.Printf("Ошибка записи в файл: %v\n", err)
					return shortURL, nil
				}
				log.Println("Успешная запись в файл", conf.FileStoragePath)
				defer prod.Close()
				return shortURL, nil
			}
		} else {
			return shortURL, nil
		}
	} else {
		shortURL, _ := uc.storage.Read(originalURL)
		return shortURL, err
	}
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

func (uc *ShortURLUseCase) Write(originalURL, shortURL string) {
	uc.storage.Create(shortURL, originalURL)
}
