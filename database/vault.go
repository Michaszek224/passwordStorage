package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
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
