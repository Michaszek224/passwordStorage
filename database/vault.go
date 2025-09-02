package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type Site struct {
	Name     string
	Password string
	Notes    string
}

func GetSiteData(db *sql.DB, userId int) []Site {
	rows, err := db.Query("SELECT site, password, notes FROM vault WHERE user_id = ?", userId)
	if err != nil {
		log.Fatalf("Error quering data from sql table: %v", err)
	}
	defer rows.Close()

	var sites []Site
	for rows.Next() {
		newSite := Site{}

		err = rows.Scan(&newSite.Name, &newSite.Password, &newSite.Notes)
		if err != nil {
			log.Fatalf("Error scanning data: %v", err)
		}
		sites = append(sites, newSite)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error after interacting with rows: %v", err)
	}

	return sites
}

func InsertSiteData(db *sql.DB, userId int, username, password, notes string) error {
	if username == "" || password == "" {
		return errors.New("site and password cannot be empty")
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO vault(user_id, site, password, notes) VALUES (?,?,?,?)", userId, username, hashedPassword, notes)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("site already taken")
		}
		return err
	}
	return nil
}
