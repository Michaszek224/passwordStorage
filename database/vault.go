package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Site struct {
	ID       int
	Name     string
	Password string
	Notes    string
}

func GetSiteData(db *sql.DB, userId int) ([]Site, error) {
	rows, err := db.Query("SELECT id, site, password, notes FROM vault WHERE user_id = ?", userId)
	if err != nil {
		log.Fatalf("Error quering data from sql table: %v", err)
	}
	defer rows.Close()

	var sites []Site
	for rows.Next() {
		newSite := Site{}

		err = rows.Scan(&newSite.ID, &newSite.Name, &newSite.Password, &newSite.Notes)
		if err != nil {
			log.Fatalf("Error scanning data: %v", err)
		}
		newSite.Password, err = passwordDecrypt(newSite.Password)
		if err != nil {
			return nil, err
		}
		sites = append(sites, newSite)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return sites, nil
}

func InsertSiteData(db *sql.DB, userId int, username, password, notes string) error {
	if username == "" || password == "" {
		return errors.New("site and password cannot be empty")
	}
	hashedPassword, err := passwordEncrypt(password)
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

func DeleteData(db *sql.DB, userId int, id string) error {
	_, err := db.Exec("DELETE FROM vault WHERE user_id = ? AND id = ?", userId, id)
	if err != nil {
		return err
	}

	return nil
}

func EditData(db *sql.DB, userId int, id, password, site, notes string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM vault WHERE user_id = ? AND site = ? AND id != ?)", userId, site, id).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("Site already exisists")
	}

	sqlSets := []string{}
	args := []interface{}{}

	if site != "" {
		sqlSets = append(sqlSets, "site = ?")
		args = append(args, site)
	}

	if password != "" {
		passwordEcrypted, err := passwordEncrypt(password)
		if err != nil {
			return err
		}
		sqlSets = append(sqlSets, "password = ?")
		args = append(args, passwordEcrypted)
	}

	if notes != "" {
		sqlSets = append(sqlSets, "notes = ?")
		args = append(args, notes)
	}

	if len(sqlSets) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf("UPDATE vault SET %s WHERE user_id = ? AND id = ?", strings.Join(sqlSets, ", "))
	args = append(args, userId, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func GetPassword(db *sql.DB, userId int, siteId string) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM vault WHERE user_Id = ? AND id = ?", userId, siteId).Scan(&password)
	if err != nil {
		return "", err
	}

	decryptedPassword, err := passwordDecrypt(password)
	if err != nil {
		return "", err
	}

	return decryptedPassword, nil
}

func passwordEncrypt(password string) (string, error) {
	key := os.Getenv("VAULT_KEY")

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(password), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func passwordDecrypt(password string) (string, error) {
	key := os.Getenv("VAULT_KEY")

	data, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
