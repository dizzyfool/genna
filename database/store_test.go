package database

import (
	"testing"

	"github.com/dizzyfool/genna/model_old"

	"go.uber.org/zap"
)

func TestLive(t *testing.T) {

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	logger, _ := config.Build()
	url := `postgres://genna:genna@localhost:5432/genna?sslmode=disable`
	db, _ := NewDatabase(url, logger)

	store := NewStore(db)

	createTestTable(store)

	var schemas []string
	schemas = append(schemas, model_old.PublicSchema)
	rows, _ := store.queryTables(schemas)

	t.Run("Should be 5 rows", func(t *testing.T) {
		got := len(rows)
		if got != 5 {
			t.Errorf("There should be 5 rows but I see %v", got)
		}
	})

	t.Run("Should not duplicate columns", func(t *testing.T) {
		if got := isThereDuplicates(rows); got {
			t.Error("Duplicated columns detected")
		}
	})

}

// Create table with complex constraints for testing
func createTestTable(store *Store) {
	_, _ = store.db.Model().Exec(`
	drop table if exists contracts;
	drop table if exists employees;

	CREATE TABLE employees (
		id serial NOT null,
		CONSTRAINT employees_pk PRIMARY KEY (id)
	);

	CREATE TABLE contracts (
		employee_id int4 NOT NULL,
		email varchar NOT NULL,
		id1 int NOT NULL,
		id2 varchar NOT NULL,
		CONSTRAINT contracts_un1 CHECK (id1 > 0),
		CONSTRAINT contracts_un2 UNIQUE (email),
		CONSTRAINT contracts_un3 UNIQUE (email, employee_id),
		CONSTRAINT contracts_un4 UNIQUE (id1, id2),
		CONSTRAINT contracts_employees_fk FOREIGN KEY (employee_id) REFERENCES employees(id)
	);
	`)
}

// Check if there are duplicated columns
func isThereDuplicates(rows []columnRow) bool {
	encountered := map[string]bool{}
	for _, v := range rows {
		if encountered[v.TblName+"|||"+v.ColumnName] {
			return true
		}

		encountered[v.TblName+"|||"+v.ColumnName] = true
	}
	return false
}
