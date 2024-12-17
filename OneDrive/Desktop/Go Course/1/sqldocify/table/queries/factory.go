package queries

import (
	"fmt"
	"sqldocify/configs"
)

func GetQueryGenerator() (QueryGenerator, error) {
	dbType := configs.GetDBType()
	switch dbType {
	case "mysql":
		return &MySQLQueryGenerator{}, nil
	// Add more cases for other database types if needed
	// case "postgresql":
	//     return &queries.PostgreSQLQueryGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
