package storage

type UserURLs struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

// Repository представляет интерфейс для работы с хранилищем данных.
type Repository interface {
	Create(short, shortURL, originURL, userID string) error
	Read(short string) (string, string, error)
	ReadOriginal(originalURL string) (string, error)
	ReadUrlsUserID(userID string) ([]UserURLs, error)
}
