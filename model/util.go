package model

import "strings"

const (
	PublicSchema   = "public"
	DefaultPackage = "model"
)

func Split(input string) (string, string) {
	d := strings.Split(input, ".")
	if len(d) < 2 {
		return PublicSchema, input
	}

	return d[0], d[1]
}

func Join(schema, table string) string {
	return schema + "." + table
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

// DiscloseSchemas discloses "*" in schemas
func DiscloseSchemas(userInput []string, tables []Table) []string {
	index := map[string]struct{}{}

	for _, t := range userInput {
		schema, table := Split(t)

		if table == "*" {
			for _, m := range tables {
				if m.Schema == schema {
					index[Join(m.Schema, m.Name)] = struct{}{}
				}
			}
		} else {
			for _, m := range tables {
				if m.Schema == schema && m.Name == table {
					index[Join(m.Schema, m.Name)] = struct{}{}
				}
			}
		}
	}

	result := make([]string, 0, len(index))
	for key := range index {
		result = append(result, key)
	}

	return result
}

// FollowFKs resolves foreign keys & adds tables in to generate list
func FollowFKs(disclosed []string, tables []Table) []string {
	iTables := map[string]struct{}{}
	for _, d := range disclosed {
		iTables[d] = struct{}{}
	}

	iModels := map[string]int{}
	for i, t := range tables {
		iModels[Join(t.Schema, t.Name)] = i
	}

	count := len(disclosed)
	for i := 0; i < count; i++ {
		table := disclosed[i]
		if m, ok := iModels[table]; ok {
			for _, r := range tables[m].Relations {
				relTable := Join(r.TargetSchema, r.TargetTable)
				_, hasModel := iModels[relTable]
				_, exists := iTables[relTable]

				if hasModel && !exists {
					iTables[relTable] = struct{}{}
					disclosed = append(disclosed, relTable)
					count++
				}
			}
		}
	}

	return disclosed
}
