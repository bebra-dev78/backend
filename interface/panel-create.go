package api

import (
	Logus "backend/logging"
	"database/sql"
	"encoding/json"
	"net/http"
)

func PanelsCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		json.NewEncoder(w).Encode(400)
		return
	}

	result, err := db.Exec("INSERT INTO panels (uid, title) VALUES (?, ?)", uid, title)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	lastInsertId, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(lastInsertId)
}
