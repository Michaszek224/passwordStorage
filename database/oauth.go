package database

import "database/sql"

func FindOrCreateOAuthUser(provider, providerID, email, username string, db *sql.DB) (int64, error) {
	var id int64
	err := db.QueryRow(`
        SELECT id FROM user WHERE provider = ? AND provider_id = ?
    `, provider, providerID).Scan(&id)

	if err == sql.ErrNoRows {
		res, insertErr := db.Exec(`
            INSERT INTO user (username, email, provider, provider_id)
            VALUES (?, ?, ?, ?)
        `, username, email, provider, providerID)
		if insertErr != nil {
			return 0, insertErr
		}
		return res.LastInsertId()
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
