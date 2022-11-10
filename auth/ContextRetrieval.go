package auth

import (
	"context"
	"database/sql"
	"log"
)

// Extracts sql.DB pointer from context object.
func GetDBFromContext(c context.Context) *sql.DB {
	db, ok := c.Value("db").(*sql.DB)
	if !ok {
		log.Panic("No database context found")
	}

	return db
}
