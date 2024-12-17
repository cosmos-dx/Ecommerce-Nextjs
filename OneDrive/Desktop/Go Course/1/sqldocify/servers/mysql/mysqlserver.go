package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"sqldocify/validators"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLServer struct {
	DB        *sql.DB
	Validator *validators.ServerValidator
}

func (m *MySQLServer) Connect(config string) (*sql.DB, error) {
	if valid, err := m.Validator.ValidateConnectionPath(config); !valid || err != nil {
		return nil, err
	}
	parts := strings.Split(config, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid config format: %s", config)
	}
	connectionString := parts[0] + "/"
	dbName := parts[1]
	dbName = strings.TrimSuffix(dbName, "?params")

	fmt.Println("config------", connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists string
	query := fmt.Sprintf("SHOW DATABASES LIKE '%s'", dbName)
	err = db.QueryRow(query).Scan(&exists)

	if err == sql.ErrNoRows {
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		_, err = db.Exec(createDBQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	db, err = sql.Open("mysql", config)
	if err != nil {
		return nil, err
	}

	m.DB = db
	return db, nil
}

func (m *MySQLServer) Close() error {
	if m.DB != nil {
		return m.DB.Close()
	}
	return errors.New("no active MySQL connection to close")
}

func (m *MySQLServer) GetDB() *sql.DB {
	return m.DB
}
