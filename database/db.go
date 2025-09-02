package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DbInit() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
			);
		`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vault(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			site TEXT NOT NULL,
			password TEXT NOT NULL,
			notes TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id),
			UNIQUE(user_id, site)
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
