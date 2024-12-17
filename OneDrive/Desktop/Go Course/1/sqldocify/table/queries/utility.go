package queries

import (
	"database/sql"
	"fmt"
	"strings"
)

func SanitizeValues(values []interface{}) string {
	var sanitizedValues []string
	for _, value := range values {
		switch v := value.(type) {
		case string:
			// Escape single quotes by doubling them
			sanitizedValues = append(sanitizedValues, fmt.Sprintf("'%s'", escapeSQLString(v)))
		case int, int32, int64, float32, float64:
			// Numbers and floats are used directly
			sanitizedValues = append(sanitizedValues, fmt.Sprintf("%v", v))
		case bool:
			// Booleans are converted to 1 or 0
			if v {
				sanitizedValues = append(sanitizedValues, "1")
			} else {
				sanitizedValues = append(sanitizedValues, "0")
			}
		case nil:
			// Handle null values
			sanitizedValues = append(sanitizedValues, "NULL")
		default:
			// For other types, try to convert them to string, but may need further handling
			sanitizedValues = append(sanitizedValues, fmt.Sprintf("'%v'", v))
		}
	}
	return strings.Join(sanitizedValues, ", ")
}

func escapeSQLString(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}

func FormatColumns(columns []string) string {
	return strings.Join(columns, ", ")
}

func SanitizeValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", escapeSQLString(v))
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
func nilIfEmpty(nullString sql.NullString) interface{} {
	if nullString.Valid {
		return nullString.String
	}
	return nil
}
