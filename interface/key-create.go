package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	Logus "backend/logging"
)

type KEY_MODEL struct {
	Api      string  `json:"api"`
	Secret   string  `json:"secret"`
	Title    string  `json:"title"`
	Phrase   *string `json:"phrase"`
	Exchange int     `json:"exchange"`
}

func KeyCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	var key KEY_MODEL

	if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
		json.NewEncoder(w).Encode(400)
		return
	}

	result, err := db.Exec("INSERT INTO `keys` (uid, api, secret, title, exchange, phrase) VALUES (?, ?, ?, ?, ?, ?)", uid, key.Api, key.Secret, key.Title, key.Exchange, key.Phrase)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	lastInsertId, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(lastInsertId)
}
