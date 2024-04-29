package storage

// Repository представляет интерфейс для работы с хранилищем данных.
type Repository interface {
	Create(url string, shortURL string) string
	Read(shortURL string) (string, error)
	Delete(shortURL string) error
}
