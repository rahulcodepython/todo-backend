package database

import "database/sql"

func createTablesInitialLoad(db *sql.DB) error {
	quires := []string{createExtensionQuery, createUserTableQuery, createJWTTableQuery, createTodoTableQuery}

	for _, query := range quires {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
