package queries

import (
	"database/sql"
	"sqldocify/configs"
)

type QueryGenerator interface { // Get schema of a table
	GenerateCreateTableQuery(nm string, schema map[string]configs.FieldSchema) string                                                                                // Single create table
	GenerateGetSchemaQuery(db *sql.DB, nm string) (map[string]configs.FieldSchema, error)                                                                            // Get schema of a table
	GenerateGetAllTablesQuery(db *sql.DB) ([]string, error)                                                                                                          // Get all tables
	GenerateTableExistsQuery(table string) string                                                                                                                    // Single table exists
	GenerateInsertQuery(table string, columns []string, values []interface{}) string                                                                                 // Single insert
	GenerateMultipleInsertQuery(table string, columns []string, values [][]interface{}) string                                                                       // Bulk insert
	GenerateUpdateQuery(table string, updates map[string]interface{}, condition string) string                                                                       // Single update
	GenerateDeleteQuery(table string, condition string) string                                                                                                       // Single delete
	GenerateMultipleDeleteQuery(table string, conditions []string) string                                                                                            // Bulk delete
	GenerateSelectQuery(table string, columns []string, condition string, orderBy string, limit int, offset int) string                                              // Single select
	GenerateBatchInsertQuery(table string, columns []string, batchValues [][]interface{}) string                                                                     // Batch insert
	GenerateUpsertQuery(table string, columns []string, values []interface{}, conflictColumns []string, updates map[string]interface{}) string                       // Single upsert
	GenerateJoinQuery(mainTable string, joinType string, joinTable string, onCondition string, columns []string, condition string, orderBy string, limit int) string // Single join
	GenerateCountQuery(table string, condition string) string                                                                                                        // Single count
	GenerateExistsQuery(table string, condition string) string                                                                                                       // Single exists
	GenerateTransactionQuery(queries []string) string                                                                                                                // Single transaction
	GenerateAggregationQuery(table string, columns []string, aggregations map[string]string, condition string, groupBy []string, orderBy string) string              // Single aggregation
	BuildConditionQuery(conditions map[string]interface{}, logicalOperator string) string                                                                            // Build condition query
	GeneratePaginationQuery(table string, columns []string, condition string, orderBy string, page int, pageSize int) string                                         // Single pagination
	GenerateCreateIndexQuery(indexName string, table string, columns []string, unique bool) string                                                                   // Single create index
	GenerateDropIndexQuery(indexName string) string                                                                                                                  // Single drop index
	GenerateAddColumnQuery(table string, columnName string, columnType string, defaultValue interface{}) string                                                      // Single add column
	GenerateModifyColumnQuery(table string, columnName string, columnType string, nullable bool) string                                                              // Single modify column
	GenerateDropColumnQuery(table string, columnName string) string                                                                                                  // Single drop column
	GenerateAddForeignKeyQuery(table string, columnName string, referencedTable string, referencedColumn string, onDelete string, onUpdate string) string            // Single add foreign key
	GenerateDropForeignKeyQuery(table string, foreignKeyName string) string                                                                                          // Single drop foreign key
	SanitizeValue(value interface{}) string                                                                                                                          // Sanitize value
	FormatColumns(columns []string) string                                                                                                                           // Format columns
	EscapeIdentifier(identifier string) string                                                                                                                       // Escape identifier
}
