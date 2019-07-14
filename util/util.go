package util

import (
	"fmt"
	"strings"
)

const (
	// PublicSchema is a default postgresql schema
	PublicSchema = "public"

	// DefaultPackage is a default package name
	DefaultPackage = "model"

	// DefaultAlias is a default alias for model
	DefaultAlias = "t"
)

// Split splits full table name in schema and table name
func Split(s string) (string, string) {
	d := strings.Split(s, ".")
	if len(d) < 2 {
		return PublicSchema, s
	}

	return d[0], d[1]
}

// Join joins table name and schema to full name
func Join(schema, table string) string {
	return schema + "." + table
}

// JoinF joins table name and schema to full name filtering public
func JoinF(schema, table string) string {
	if schema == PublicSchema {
		return table
	}

	return Join(schema, table)
}

// Quoted quotes entity name if needed
func Quoted(fullName string, escape bool) string {
	if !HasUpper(fullName) {
		return fullName
	}

	pattern := `"%s"`
	if escape {
		pattern = `\"%s\"`
	}

	d := strings.Split(fullName, ".")
	if len(d) < 2 {
		return fmt.Sprintf(pattern, fullName)
	}

	return Join(fmt.Sprintf(pattern, d[0]), fmt.Sprintf(pattern, d[1]))
}

// Schemas get schemas from table names
func Schemas(tables []string) (schemas []string) {
	index := map[string]struct{}{}
	for _, t := range tables {
		schema, _ := Split(t)
		if _, ok := index[schema]; !ok {
			index[schema] = struct{}{}
			schemas = append(schemas, schema)
		}
	}

	return
}
