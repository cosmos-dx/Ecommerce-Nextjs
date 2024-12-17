package configs

import (
	"database/sql"
	"sync"
)

var dbType string

func SetDBType(db string) {
	dbType = db
}

func GetDBType() string {
	return dbType
}

type DBServer interface {
	Connect(config string) (*sql.DB, error)
	Close() error
	GetDB() *sql.DB
}

type Database struct {
	DBServer DBServer
}

func (d *Database) DB() *sql.DB {
	return d.DBServer.GetDB()
}

func (d *Database) Close() error {
	return d.DBServer.Close()
}

type MetaTableDetails struct {
	Schema    map[string]FieldSchema `json:"schema"`
	Timestamp string                 `json:"timestamp"`
	Details   string                 `json:"details"`
}

type MetaTableList struct {
	mu             sync.Mutex
	ExistingTables map[string]MetaTableDetails `json:"existing_tables"`
}

var (
	instance *MetaTableList
	once     sync.Once
)

type FieldSchema struct {
	Type string `json:"Type"`
	Null string `json:"Null"`
	Key  string `json:"Key"`
	// Default interface{} `json:"Default"`
	Extra string `json:"Extra"`
}
