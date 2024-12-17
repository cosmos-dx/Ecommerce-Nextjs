package table

import "sqldocify/configs"

type UpdateDiffs struct {
	FieldName string
	OldValue  interface{}
	NewValue  interface{}
}

type ITableSpec interface {
	GetMetaDataSchema() (interface{}, error)
	TableExists(nm string, db *configs.Database) bool
	CreateTable(db *configs.Database, nm string, schema map[string]configs.FieldSchema) error
	Insert(dt interface{}, db *configs.Database) error
	Update(dt interface{}, onUpdate func() error, db *configs.Database) (UpdateDiffs, error)
	Delete(condition interface{}, db *configs.Database) error
	Fetch(condition interface{}, result interface{}, db *configs.Database) error
	BeginTransaction(db *configs.Database) error
	CommitTransaction(db *configs.Database) error
	RollbackTransaction(db *configs.Database) error
	BatchInsert(dts []interface{}, db *configs.Database) error
	BatchUpdate(dts []interface{}, onUpdate func() error, db *configs.Database) ([]UpdateDiffs, error)
	BatchDelete(conditions []interface{}, db *configs.Database) error
}
