package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

func GetUsers(db *sql.DB) []User {
	rows, err := db.Query("SELECT username, password FROM User")
	if err != nil {
		log.Fatalf("Error quering data from sql table: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		newUser := User{}

		err := rows.Scan(&newUser.Username, &newUser.Password)
		if err != nil {
			log.Fatalf("Error scaning data: %v", err)
		}
		users = append(users, newUser)

	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error after interacting with rows: %v", err)
	}
	return users
}

func InsertUser(username, password string, db *sql.DB) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO user(username, password) VALUES (?,?)`, username, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("username already taken")
		}
		return err
	}
	return err
}

func AuthenicateUser(username, password string, db *sql.DB) error {
	var hashedPassword string

	err := db.QueryRow(
		`SELECT password FROM user WHERE username = ?`, username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return errors.New("User not found")
	}

	if err != nil {
		return err
	}

	if !CheckHashPassowrd(password, hashedPassword) {
		return errors.New("invalid password")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckHashPassowrd(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
