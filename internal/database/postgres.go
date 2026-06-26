package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

// NewConnection creates a new database connection.
func NewConnection(connectionString string) (*sql.DB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("could not open database connection: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("could not ping database: %w", err)
    }

    log.Println("Database connection established")
    return db, nil
}
