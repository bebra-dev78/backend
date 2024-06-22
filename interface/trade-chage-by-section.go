package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func TradeChangeBySection(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	section := r.URL.Query().Get("section")
	id := r.URL.Query().Get("id")
	if section == "" || id == "" {
		return
	}

	switch section {
	case "rating":
		value := r.URL.Query().Get("value")
		if value == "" {
			return
		}

		db.Exec("UPDATE trades SET rating = ? WHERE id = ?", value, id)
	case "tags":
		var tags []string

		if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
			return
		}

		data, _ := json.Marshal(tags)

		db.Exec("UPDATE trades SET tags = ? WHERE id = ?", string(data), id)
	default:
	}
}
