package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func TradesDeleteByIds(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	json.NewEncoder(w).Encode(200)
}
