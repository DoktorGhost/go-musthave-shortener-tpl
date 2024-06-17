package postgres

import (
	"database/sql"
	"errors"
	"sync"
)

type PostgresStorage struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewPostgresRepository создает новый экземпляр PostgresRepository.
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Создание таблицы, если она не существует
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        url TEXT NOT NULL UNIQUE,
        short_url TEXT NOT NULL UNIQUE
    );
    `
	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (r *PostgresStorage) Read(originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var shortURL string
	err := r.db.QueryRow("SELECT url FROM urls WHERE short_url = $1", originalURL).Scan(&shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("url not found")
		}
		return "", err
	}
	return shortURL, nil
}

// возвращает сокращенный URL и TRUE, если его еще не было в БД
func (r *PostgresStorage) Create(shortURL string, url string) (string, bool) {
	//читаем из БД по оригинальной ссылке
	val, err := r.Read(url)
	//если ошибка есть, то записываем значение в БД
	if err != nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		query := "INSERT INTO urls (url, short_url) VALUES ($1, $2) RETURNING short_url"
		var returnedShortURL string
		err := r.db.QueryRow(query, url, shortURL).Scan(&returnedShortURL)
		if err != nil {
			return "", false
		}
		query = "INSERT INTO urls (url, short_url) VALUES ($1, $2)"
		_, err = r.db.Exec(query, shortURL, url)
		if err != nil {
			return "", false
		}
		return returnedShortURL, true
	}
	return val, false
}

func (r *PostgresStorage) Delete(shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM urls WHERE short_url = $1 AND url = $1"
	_, err := r.db.Exec(query, shortURL)
	if err != nil {
		return err
	}
	return nil
}
