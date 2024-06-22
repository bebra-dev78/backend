package api

import (
	Logus "backend/logging"
	"database/sql"
	"net/http"
)

func KeyChangeTitle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		return
	}

	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	if id == "" || title == "" {
		return
	}

	if _, err := db.Exec("UPDATE `keys` SET title = ? WHERE id = ?", title, id); err != nil {
		Logus.Logus.Error(err)
		return
	}
}
