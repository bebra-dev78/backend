package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"unicode/utf8"

	_ "backend/logging"
)

func UserChangeBySocial(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	name := r.URL.Query().Get("name")
	public := r.URL.Query().Get("public")
	surname := r.URL.Query().Get("surname")

	if utf8.RuneCountInString(name) < 3 || utf8.RuneCountInString(name) > 13 || utf8.RuneCountInString(public) > 20 {
		json.NewEncoder(w).Encode(405)
		return
	}

	if public == "" {
		if _, err := db.Exec("UPDATE users SET name = ?, surname = ?, public = NULL WHERE id = ?", name, surname, uid); err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}
	} else {
		var existingUID string
		if err := db.QueryRow("SELECT id FROM users WHERE public = ?", public).Scan(&existingUID); err != nil && err != sql.ErrNoRows {
			json.NewEncoder(w).Encode(500)

			return
		} else if err == sql.ErrNoRows || existingUID == uid {
			_, err := db.Exec("UPDATE users SET name = ?, surname = ?, public = ? WHERE id = ?", name, surname, public, uid)
			if err != nil {
				json.NewEncoder(w).Encode(nil)
				return
			}
		} else {
			json.NewEncoder(w).Encode(409)
			return
		}
	}

	json.NewEncoder(w).Encode(200)
}
