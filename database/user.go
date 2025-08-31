package database

import (
	"database/sql"
	"log"
)

type User struct {
	Nickname string
	Password string
}

func GetUsers(db *sql.DB) []User {
	rows, err := db.Query("SELECT nickname, password FROM User")
	if err != nil {
		log.Fatalf("Error quering data from sql table: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var nickname string
		var password string

		err := rows.Scan(&nickname, &password)
		if err != nil {
			log.Fatalf("Error scaning data: %v", err)
		}
		users = append(users, User{Nickname: nickname, Password: password})

	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error after interacting with rows: %v", err)
	}
	return users
}
