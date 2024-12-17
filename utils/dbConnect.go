package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnect() (*sql.DB, error) {
	dsn := "root:root@tcp(localhost:3306)/stockscreener"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error verifying database connection: %v", err)
	}
	return db, nil
}
