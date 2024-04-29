package usecase

type ShortUrlUseCaseInterface interface {
	CreateShortUrl(originalURL string) string
	GetOriginalUrl(shortURL string) (string, error)
	DeleteUrl(shortURL string) error
}
