package storage

// Repository представляет интерфейс для работы с хранилищем данных.
type Repository interface {
	Create(url string, shortURL string) error
	Read(shortURL string) (string, error)
	Update(shortURL string, newURL string) error
	Delete(shortURL string) error
}
