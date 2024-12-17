package sqlite

import (
	"database/sql"
	"errors"
	"sqldocify/validators"
	// _ "github.com/mattn/go-sqlite3"
)

type SQLiteServer struct {
	DB        *sql.DB
	Validator *validators.ServerValidator
}

func (s *SQLiteServer) Connect(config string) (*sql.DB, error) {
	if valid, err := s.Validator.ValidateConnectionPath(config); !valid || err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", config)
	if err != nil {
		return nil, err
	}
	s.DB = db
	return db, nil
}

func (s *SQLiteServer) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return errors.New("no active SQLite connection to close")
}

func (s *SQLiteServer) GetDB() *sql.DB {
	return s.DB
}
