package table

import (
	"database/sql"
	"errors"
	"log"
	"sqldocify/configs"
	"sqldocify/table/queries"
)

var QGType queries.QueryGenerator

type TableSpec struct {
	TableName string
	QGType    queries.QueryGenerator
}

func AddSelectedDB() (*TableSpec, error) {
	queryGenerator, err := queries.GetQueryGenerator() // Load the query generator once
	if err != nil {
		return nil, err
	}
	return &TableSpec{
		QGType: queryGenerator,
	}, nil
}

func (t *TableSpec) GetSelectedDB() queries.QueryGenerator {
	return t.QGType
}
func (t *TableSpec) GetAllTablesList(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, errors.New("no active database connection")
	}
	return t.QGType.GenerateGetAllTablesQuery(db)
}

func (t *TableSpec) GetMetaDataSchema(db *sql.DB, tname string) (map[string]configs.FieldSchema, error) {
	return t.QGType.GenerateGetSchemaQuery(db, tname)
	// return nil, nil
}
func (t *TableSpec) TableExists(nm string, db *sql.DB) bool {
	// exsisting_tables := configs.GetTableListInstance()
	// if()
	//check from metatables
	tableExistsQuery := t.QGType.GenerateTableExistsQuery(nm)
	if db == nil {
		errors.New("No active database connection.")
		return false
	}
	db.Exec(tableExistsQuery)
	return true
}
func (t *TableSpec) CreateTable(db *sql.DB, nm string, schema map[string]configs.FieldSchema) error {
	createQuery := t.QGType.GenerateCreateTableQuery(nm, schema)
	if db == nil {
		return errors.New("no active database connection")
	}
	_, err := db.Exec(createQuery)
	return err
}

func (t *TableSpec) Insert(dt interface{}, db *sql.DB) error {
	columns := []string{"column1", "column2"}
	values := []interface{}{dt, dt}

	insertQuery := t.QGType.GenerateInsertQuery(t.TableName, columns, values)
	log.Printf("Insert Query: %s", insertQuery)
	if db == nil {
		return errors.New("no active database connection")
	}
	return nil
}
