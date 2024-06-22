package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func KeyDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		return
	}

	db.Exec("DELETE FROM `keys` WHERE id = ?", id)
}
