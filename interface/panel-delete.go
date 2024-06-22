package api

import (
	"database/sql"
	"net/http"
)

func PanelsDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		return
	}

	db.Exec("DELETE FROM panels WHERE id = ?", id)
}
