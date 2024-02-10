package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	op := "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	); 
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, err
}
func (s *Storage) SaveURL(urlToSave string, alias string) error {
	op := "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		// TODO refactor this
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	op := "storage.sqlite.GetURL"

	var url string
	row := s.db.QueryRow("SELECT url FROM url WHERE alias = $1 LIMIT 1", alias)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	err := row.Scan(&url)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, err
}
