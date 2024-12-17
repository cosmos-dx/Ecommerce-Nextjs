package queries

import (
	"database/sql"
	"fmt"
	"sqldocify/configs"
	"strings"
)

type MySQLQueryGenerator struct{}

func (m *MySQLQueryGenerator) GenerateGetSchemaQuery(db *sql.DB, tablename string) (map[string]configs.FieldSchema, error) {
	query := "DESC " + tablename + ";"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schema := make(map[string]configs.FieldSchema)

	for rows.Next() {
		var field string
		var fieldType string
		var isNull string
		var key string
		var defaultValue sql.NullString
		var extra string

		err := rows.Scan(&field, &fieldType, &isNull, &key, &defaultValue, &extra)
		if err != nil {
			return nil, err
		}

		schema[field] = configs.FieldSchema{
			Type: fieldType,
			Null: isNull,
			Key:  key,
			// Default: nilIfEmpty(defaultValue),
			Extra: extra,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schema, nil
}

func (m *MySQLQueryGenerator) GenerateGetAllTablesQuery(db *sql.DB) ([]string, error) {
	query := "SHOW TABLES;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	tables := []string{}
	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return tables, nil
}

func (m *MySQLQueryGenerator) GenerateCreateTableQuery(nm string, schema map[string]configs.FieldSchema) string {
	var columnStrings []string
	for column, fieldSchema := range schema {
		if column == "created_at" || column == "updated_at" {
			continue
		}
		colDef := fmt.Sprintf("%s %s", column, fieldSchema.Type)

		if fieldSchema.Null != "" {
			colDef += " " + "NOT NULL"
		}
		if fieldSchema.Key == "UNI" {
			colDef += " " + "UNIQUE"
		}
		if fieldSchema.Key == "PRI" {
			colDef += " " + "PRIMARY KEY"
		}
		if fieldSchema.Extra != "" {
			colDef += " " + fieldSchema.Extra
		}
		columnStrings = append(columnStrings, colDef)
	}

	return fmt.Sprintf("CREATE TABLE %s (%s);", nm, strings.Join(columnStrings, ", "))
}

func (m *MySQLQueryGenerator) GenerateTableExistsQuery(table string) string {
	return fmt.Sprintf("SHOW TABLES LIKE '%s';", table)
}

func (m *MySQLQueryGenerator) GenerateInsertQuery(table string, columns []string, values []interface{}) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, FormatColumns(columns), SanitizeValues(values))
}

func (m *MySQLQueryGenerator) GenerateMultipleInsertQuery(table string, columns []string, values [][]interface{}) string {
	var valueStrings []string
	for _, val := range values {
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", SanitizeValues(val)))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, FormatColumns(columns), strings.Join(valueStrings, ", "))
}

func (m *MySQLQueryGenerator) GenerateUpdateQuery(table string, updates map[string]interface{}, condition string) string {
	var setClauses []string
	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = %s", column, SanitizeValue(value)))
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s;", table, strings.Join(setClauses, ", "), condition)
}
func (m *MySQLQueryGenerator) GenerateDeleteQuery(table string, condition string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s;", table, condition)
}

func (m *MySQLQueryGenerator) GenerateMultipleDeleteQuery(table string, conditions []string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s;", table, strings.Join(conditions, " OR "))
}
func (m *MySQLQueryGenerator) GenerateSelectQuery(table string, columns []string, condition string, orderBy string, limit int, offset int) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d OFFSET %d;", FormatColumns(columns), table, condition, orderBy, limit, offset)
}

func (m *MySQLQueryGenerator) GenerateBatchInsertQuery(table string, columns []string, batchValues [][]interface{}) string {
	var valueStrings []string
	for _, val := range batchValues {
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", SanitizeValues(val)))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", table, FormatColumns(columns), strings.Join(valueStrings, ", "))
}

func (m *MySQLQueryGenerator) GenerateUpsertQuery(table string, columns []string, values []interface{}, conflictColumns []string, updates map[string]interface{}) string {
	var setClauses []string
	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = %s", column, SanitizeValue(value)))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s;", table, FormatColumns(columns), SanitizeValues(values), strings.Join(setClauses, ", "))
}

func (m *MySQLQueryGenerator) GenerateJoinQuery(mainTable string, joinType string, joinTable string, onCondition string, columns []string, condition string, orderBy string, limit int) string {
	return fmt.Sprintf("SELECT %s FROM %s %s JOIN %s ON %s WHERE %s ORDER BY %s LIMIT %d;", FormatColumns(columns), mainTable, joinType, joinTable, onCondition, condition, orderBy, limit)
}

func (m *MySQLQueryGenerator) GenerateCountQuery(table string, condition string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s;", table, condition)
}

func (m *MySQLQueryGenerator) GenerateExistsQuery(table string, condition string) string {
	return fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s);", table, condition)
}

func (m *MySQLQueryGenerator) GenerateTransactionQuery(queries []string) string {
	return fmt.Sprintf("START TRANSACTION;\n%s;\nCOMMIT;", strings.Join(queries, ";\n"))
}
func (m *MySQLQueryGenerator) GenerateAggregationQuery(table string, columns []string, aggregations map[string]string, condition string, groupBy []string, orderBy string) string {
	var aggClauses []string
	for column, aggFunc := range aggregations {
		aggClauses = append(aggClauses, fmt.Sprintf("%s(%s)", aggFunc, column))
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s GROUP BY %s ORDER BY %s;", strings.Join(aggClauses, ", "), table, condition, strings.Join(groupBy, ", "), orderBy)
}
func (m *MySQLQueryGenerator) BuildConditionQuery(conditions map[string]interface{}, logicalOperator string) string {
	var conditionClauses []string
	for column, value := range conditions {
		conditionClauses = append(conditionClauses, fmt.Sprintf("%s = %s", column, SanitizeValue(value)))
	}
	return strings.Join(conditionClauses, fmt.Sprintf(" %s ", logicalOperator))
}

func (m *MySQLQueryGenerator) GeneratePaginationQuery(table string, columns []string, condition string, orderBy string, page int, pageSize int) string {
	offset := (page - 1) * pageSize
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d OFFSET %d;", FormatColumns(columns), table, condition, orderBy, pageSize, offset)
}
func (m *MySQLQueryGenerator) GenerateCreateIndexQuery(indexName string, table string, columns []string, unique bool) string {
	uniqueClause := ""
	if unique {
		uniqueClause = "UNIQUE "
	}
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s);", uniqueClause, indexName, table, FormatColumns(columns))
}
func (m *MySQLQueryGenerator) GenerateDropIndexQuery(indexName string) string {
	return fmt.Sprintf("DROP INDEX %s;", indexName)
}
func (m *MySQLQueryGenerator) GenerateAddColumnQuery(table string, columnName string, columnType string, defaultValue interface{}) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s DEFAULT %s;", table, columnName, columnType, SanitizeValue(defaultValue))
}
func (m *MySQLQueryGenerator) GenerateModifyColumnQuery(table string, columnName string, columnType string, nullable bool) string {
	nullableClause := ""
	if nullable {
		nullableClause = "NULL"
	} else {
		nullableClause = "NOT NULL"
	}
	return fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s %s;", table, columnName, columnType, nullableClause)
}
func (m *MySQLQueryGenerator) GenerateDropColumnQuery(table string, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", table, columnName)
}

func (m *MySQLQueryGenerator) GenerateAddForeignKeyQuery(table string, columnName string, referencedTable string, referencedColumn string, onDelete string, onUpdate string) string {
	return fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT FK_%s FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE %s ON UPDATE %s;", table, columnName, columnName, referencedTable, referencedColumn, onDelete, onUpdate)
}
func (m *MySQLQueryGenerator) GenerateDropForeignKeyQuery(table string, foreignKeyName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s;", table, foreignKeyName)
}

func (m *MySQLQueryGenerator) SanitizeValue(value interface{}) string {
	return SanitizeValue(value)
}

func (m *MySQLQueryGenerator) FormatColumns(columns []string) string {
	return FormatColumns(columns)
}

func (m *MySQLQueryGenerator) EscapeIdentifier(identifier string) string {
	return fmt.Sprintf("`%s`", identifier)
}
