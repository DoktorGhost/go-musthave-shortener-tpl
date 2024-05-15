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

func (uc *ShortURLUseCase) CreateShortURL(originalURL string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", err
	}

	for i := 0; i < 10; i++ {
		//рандомная строка, будующая ссылка
		shortURL := shortener.RandomString(8)
		//проверяем, что данной строки нет в бд
		_, err := uc.storage.Read(shortURL)
		//если ошибка есть, значит данной строки нет в БД
		if err != nil {
			//записываем
			short, flags := uc.storage.Create(shortURL, originalURL)
			//запись в файл, если флаг true

			if flags {
				if config.FileStoragePath != "" {
					prod, err := osfile.NewProducer(config.FileStoragePath)
					if err != nil {
						log.Printf("Ошибка создания Producer: %v\n", err)
						return short, nil
					} else {
						currentTime := time.Now()
						intFromTime := currentTime.Unix()
						event := osfile.Event{
							UUID:        strconv.Itoa(int(intFromTime)),
							ShortURL:    short,
							OriginalURL: originalURL,
						}
						err = prod.WriteEvent(&event)
						if err != nil {
							log.Printf("Ошибка записи в файл: %v\n", err)
							return short, nil
						}
						log.Println("Успешная запись в файл", config.FileStoragePath)
						defer prod.Close()
					}

				}
			}

			return short, nil
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

func (uc *ShortURLUseCase) Write(originalURL, shortURL string) {
	uc.storage.Create(shortURL, originalURL)
}
