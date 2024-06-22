package api

import (
	"database/sql"
	"net/http"
)

func TradesDeleteByKey(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		return
	}

	kid := r.URL.Query().Get("kid")
	if kid == "" {
		return
	}

	db.Exec(`DELETE FROM trades WHERE kid = ?`, kid)
}
