package main

import (
	"fmt"
	"log"
	"sqldocify/configs"
	"sqldocify/servers"
	"sqldocify/table"
)

type DummyData struct {
	tableName string
	columns   []string
	values    []interface{}
}

type FieldSchema struct {
	Type  string `json:"Type"`
	Null  string `json:"Null"`
	Key   string `json:"Key"`
	Extra string `json:"Extra"`
}

var userTableSchema = map[string]configs.FieldSchema{
	"email": {
		Type:  "varchar(255)",
		Null:  "NO",
		Key:   "UNI",
		Extra: "",
	},
	"id": {
		Type:  "int",
		Null:  "NO",
		Key:   "PRI",
		Extra: "auto_increment",
	},
	"password": {
		Type:  "varchar(255)",
		Null:  "NO",
		Key:   "",
		Extra: "",
	},
	"status": {
		Type:  "enum('active','inactive','banned')",
		Null:  "NO",
		Key:   "",
		Extra: "",
	},
	"username": {
		Type:  "varchar(255)",
		Null:  "NO",
		Key:   "UNI",
		Extra: "",
	},
}

func main() {
	dbType := "mysql"
	config := "root:pass@tcp(127.0.0.1:3306)/gotest1"

	db, err := servers.NewDatabase(dbType, config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close the database connection: %v", err)
		}
	}()
	// .CreateTables(db)
	// table.TableExists("users", db)
	tablespec, _ := table.AddSelectedDB()
	tablespec.CreateTable(db.DB(), "abc", userTableSchema)
	fmt.Println("Database connection established successfully!")

	fmt.Println("Operations completed.")
}
