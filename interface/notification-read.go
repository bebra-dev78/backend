package api

import (
	Logus "backend/logging"
	"database/sql"
	"encoding/json"
	"net/http"
)

func NotificationsRead(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	_, err := db.Exec("UPDATE notifications SET `read` = true WHERE uid = ? AND `read` = false", uid)
	if err != nil {
		Logus.Logus.Error(err)
		return
	}
}
