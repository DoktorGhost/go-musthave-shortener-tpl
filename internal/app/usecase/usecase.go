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

// генерирует рандомную строку и проверяет, что в БД нет записи с этой строкой
func (uc *ShortURLUseCase) GenerateShort(originalURL string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", err
	}

	//рандомная строка, будующая ссылка

	//вставить проверку, что данной строки нет в БД
	for i := 1; i <= 3; i++ {
		short := shortener.RandomString(8)
		_, _, err := uc.storage.Read(short)
		if err != nil {
			return short, nil
		} else {
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
	}
	return "", errors.New("ошибка генерации сокрвщенного URL")
}

// запись
func (uc *ShortURLUseCase) CreateShortURL(short, shortURL, originalURL, userID string, conf *config.Config) error {

	err := uc.storage.Create(short, shortURL, originalURL, userID)

	//запись в файл
	if err == nil {
		if conf.FileStoragePath != "" {
			prod, err := osfile.NewProducer(conf.FileStoragePath)
			if err != nil {
				log.Printf("Ошибка создания Producer: %v\n", err)
				return nil
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
					return nil
				}
				log.Println("Успешная запись в файл", conf.FileStoragePath)
				defer prod.Close()
				return nil
			}
		} else {
			return nil
		}
	} else {
		return err
	}
}

// возвращает оригинальный урл по сокращенной строке "авыаыв"
func (uc *ShortURLUseCase) GetOriginalURL(short string) (string, error) {
	originalURL, _, err := uc.storage.Read(short)
	if err != nil {
		return "", err
	} else {
		return originalURL, nil
	}
}

// возвращает оригинальный урл по сокращенной строке "авыаыв"
func (uc *ShortURLUseCase) GetShortURL(short string) (string, error) {
	_, shortURL, err := uc.storage.Read(short)
	if err != nil {
		return "", err
	} else {
		return shortURL, nil
	}
}

// запись
func (uc *ShortURLUseCase) Write(short, shortURL, originalURL, userID string) error {
	return uc.storage.Create(short, shortURL, originalURL, userID)
}

// получение всех урлов от одного ИД пользователя
func (uc *ShortURLUseCase) GetUserURL(userID string) ([]storage.UserURLs, error) {
	return uc.storage.ReadUrlsUserId(userID)
}

// получение шорта по оригинальной ссылке. ошибка если не найдена запись
func (uc *ShortURLUseCase) GetShort(originalURL string) (string, error) {
	short, err := uc.storage.ReadOriginal(originalURL)

	if err != nil {
		return "", err
	} else {
		return short, nil
	}
}
