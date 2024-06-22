package api

import (
	Logus "backend/logging"
	"database/sql"
	"encoding/json"
	"net/http"
)

type Notification struct {
	Title    string `db:"title" json:"title"`
	ImageURL string `db:"image_url" json:"image_url"`
}

func NotificationsCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	var notification Notification

	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		json.NewEncoder(w).Encode(400)
		return
	}

	result, err := db.Exec("INSERT INTO notifications (uid, title, image_url) VALUES (?, ?, ?)", uid, &notification.Title, &notification.ImageURL)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	lastInsertId, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(lastInsertId)
}
