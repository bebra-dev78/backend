package resourses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	Logus "backend/logging"
)

type USER_MODEL struct {
	ID          string     `db:"id" json:"id"`
	Email       string     `db:"email" json:"email"`
	Name        string     `db:"name" json:"name"`
	Public      *string    `db:"public" json:"public"`
	Surname     *string    `db:"surname" json:"surname"`
	Private     bool       `db:"private" json:"private"`
	Convert     bool       `db:"convert" json:"convert"`
	ActivatedAt *time.Time `db:"activated_at" json:"activated_at"`
}

type KEY_MODEL struct {
	ID       uint    `db:"id" json:"id"`
	Api      string  `db:"api" json:"api"`
	Secret   string  `db:"secret" json:"secret"`
	Title    string  `db:"title" json:"title"`
	Phrase   *string `db:"phrase" json:"phrase"`
	Exchange uint8   `db:"exchange" json:"exchange"`
}

func User(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	secret := r.Header.Get("Secret")
	if secret != "bebra" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	email := r.URL.Query().Get("email")

	var user USER_MODEL

	if err := db.QueryRow("SELECT id, email, name, surname, public, private, `convert`, activated_at FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Name, &user.Surname, &user.Public, &user.Private, &user.Convert, &user.ActivatedAt); err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(nil)
		} else {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
		}
		return
	}

	rows, err := db.Query("SELECT id, api, secret, title, phrase, exchange FROM `keys` WHERE uid = ?", user.ID)
	if err != nil {
		json.NewEncoder(w).Encode(500)
		Logus.Logus.Error(err)
		return
	}

	defer rows.Close()

	var key KEY_MODEL
	var keys []KEY_MODEL

	for rows.Next() {
		if err := rows.Scan(&key.ID, &key.Api, &key.Secret, &key.Title, &key.Phrase, &key.Exchange); err != nil {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
			return
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		keys = []KEY_MODEL{}
	}

	json.NewEncoder(w).Encode([]interface{}{user, keys})
}
