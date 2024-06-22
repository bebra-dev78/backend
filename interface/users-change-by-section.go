package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"unicode/utf8"

	Logus "backend/logging"

	"golang.org/x/crypto/bcrypt"
)

func UserChangeBySection(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	section := r.URL.Query().Get("section")
	value := r.URL.Query().Get("value")
	if section == "" || value == "" {
		json.NewEncoder(w).Encode(400)
		return
	}

	switch section {
	case "private":
		db.Exec("UPDATE users SET private = ? WHERE id = ?", value, uid)
	case "convert":
		db.Exec("UPDATE users SET `convert` = ? WHERE id = ?", value, uid)
	case "password":
		password := r.URL.Query().Get("password")

		if utf8.RuneCountInString(value) < 8 || utf8.RuneCountInString(value) > 24 {
			json.NewEncoder(w).Encode(400)
			return
		}

		var userPassword string

		if err := db.QueryRow("SELECT password FROM users WHERE id = ?", uid).Scan(&userPassword); err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(nil)
			} else {
				json.NewEncoder(w).Encode(500)
				Logus.Logus.Error(err)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
			json.NewEncoder(w).Encode(403)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(value), 10)

		if _, err := db.Exec("UPDATE users SET password = ? WHERE id = ?", hash, uid); err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		json.NewEncoder(w).Encode(200)
	default:
	}
}
