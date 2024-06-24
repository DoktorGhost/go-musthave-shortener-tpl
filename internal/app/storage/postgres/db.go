package postgres

import (
	"database/sql"
	"errors"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/DoktorGhost/go-musthave-shortener-tpl/internal/app/usecase"
	"github.com/jackc/pgx/v5/pgconn"
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
        short TEXT NOT NULL UNIQUE,
        short_url TEXT NOT NULL UNIQUE,
        original_url TEXT NOT NULL UNIQUE,
        user_id TEXT NOT NULL
    );
    `
	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (r *PostgresStorage) Read(short string) (string, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var originalURL string
	var shortURL string
	err := r.db.QueryRow("SELECT original_url, short_url FROM urls WHERE short = $1", short).Scan(&originalURL, &shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errors.New("url not found")
		}
		return "", "", err
	}
	return originalURL, shortURL, nil
}

func (r *PostgresStorage) ReadOriginal(originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var short string
	err := r.db.QueryRow("SELECT short FROM urls WHERE original_url = $1", originalURL).Scan(&short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("url not found")
		}
		return "", err
	}
	return short, nil
}

func (r *PostgresStorage) ReadUrlsUserID(userID string) ([]storage.UserURLs, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var urls []storage.UserURLs
	rows, err := r.db.Query("SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url storage.UserURLs
		err := rows.Scan(&url.ShortURL, &url.LongURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return nil, errors.New("url not found")
	}

	return urls, nil

}

// возвращает сокращенный URL и TRUE, если его еще не было в БД
func (r *PostgresStorage) Create(short, shortURL, originURL, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "INSERT INTO urls (short, short_url, original_url, user_id) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, short, shortURL, originURL, userID)

	if err != nil {
		// Проверяем, является ли это ошибкой нарушения уникальности
		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			// Обработка других типов ошибок
			return err
		}

		if pgErr.Code == "23505" { // Код для ошибки нарушения уникальности
			return usecase.ErrShortURLAlreadyExists
		}

		// Обработка других ошибок
		return err
	}
	return nil
}

/*
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
*/
