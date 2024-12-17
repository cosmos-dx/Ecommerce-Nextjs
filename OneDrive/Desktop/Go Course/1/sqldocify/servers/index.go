package servers

import (
	"errors"
	"fmt"
	"log"

	"sqldocify/configs"
	"sqldocify/servers/mysql"
	"sqldocify/servers/sqlite"
	"sqldocify/table"
	"sqldocify/validators"
)

// NewDatabase initializes a new database connection and manages the active tables list.
func NewDatabase(dbtype, config string) (*configs.Database, error) {
	validator := &validators.ServerValidator{}
	var dbServer configs.DBServer

	switch dbtype {
	case "mysql":
		dbServer = &mysql.MySQLServer{Validator: validator}
		configs.SetDBType("mysql")
	case "sqlite":
		dbServer = &sqlite.SQLiteServer{Validator: validator}
		configs.SetDBType("sqlite")
	default:
		return nil, errors.New("database type not supported")
	}

	if _, err := dbServer.Connect(config); err != nil {
		return nil, err
	}
	table.AddSelectedDB()
	configs.GetMetaTableInstance() //It will load all the data in the list

	return &configs.Database{DBServer: dbServer}, nil
}

//Function for check table exists or not if not in database and if its in metadata then create the table with the available schema
/*
1. Load all the tables name from the file
2. Run for all the table names and check if they exists in db or not
3. If not then create the table with the available metadata
*/

func InitialTablesCheck(db *configs.Database) error {
	metaTables := configs.GetMetaTableInstance()
	tableSpec, err := table.AddSelectedDB()
	if err != nil {
		return err
	}
	dbtablelist, err := tableSpec.GetAllTablesList(db.DB())
	if err != nil {
		return fmt.Errorf("failed to fetch database tables: %v", err)
	}
	var metatablearray []string
	for tableName := range metaTables.ExistingTables {
		metatablearray = append(metatablearray, tableName)
	}

	tableExistsInArray := func(tableName string, tableArray []string) bool {
		for _, name := range tableArray {
			if name == tableName {
				return true
			}
		}
		return false
	}

	// 1. If any table exists in dbtablelist but not in metatablearray, update metaTables
	for _, dbTable := range dbtablelist {
		if !tableExistsInArray(dbTable, metatablearray) {
			tableschema, err := tableSpec.GetMetaDataSchema(db.DB(), dbTable)
			if err != nil {
				log.Printf("failed to get schema for table %s: %v", dbTable, err)
				continue
			}
			metatabledetails := configs.MetaTableDetails{
				Schema:    tableschema,
				Timestamp: "Timestamp",
				Details:   "Details",
			}
			metaTables.UpdateMetaTable(dbTable, metatabledetails)
			log.Printf("Table %s added to metaTables", dbTable)
		}
	}
	// 2. If any table exists in metatablearray but not in dbtablelist, create the table
	for _, metaTable := range metatablearray {
		if !tableExistsInArray(metaTable, dbtablelist) {
			tableschema := metaTables.ExistingTables[metaTable].Schema
			if err := tableSpec.CreateTable(db.DB(), metaTable, tableschema); err != nil {
				log.Printf("Failed to create table %s: %v", metaTable, err)
				continue
			}
			log.Printf("Table %s was missing and has been created", metaTable)
		}
	}

	metaTables = configs.GetMetaTableInstance()
	return nil
}
