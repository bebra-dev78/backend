package api

import (
	Logus "backend/logging"
	"database/sql"
	"encoding/json"
	"net/http"
)

func PanelsChangeWidgets(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		json.NewEncoder(w).Encode(400)
		return
	}

	var widgets []int8

	if err := json.NewDecoder(r.Body).Decode(&widgets); err != nil {
		json.NewEncoder(w).Encode(401)
		Logus.Logus.Error(err)
		return
	}

	data, _ := json.Marshal(widgets)

	if _, err := db.Exec("UPDATE panels SET widgets = ? WHERE id = ?", string(data), id); err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}
}
