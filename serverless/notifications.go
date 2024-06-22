package resourses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	Logus "backend/logging"
)

type NOTIFICATION_MODEL struct {
	ID        uint      `db:"id" json:"id"`
	Read      bool      `db:"read" json:"read"`
	Title     string    `db:"title" json:"title"`
	ImageURL  string    `db:"image_url" json:"image_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func Notifications(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	rows, err := db.Query("SELECT id, `read`, title, image_url, created_at FROM notifications WHERE uid = ?", uid)
	if err != nil {
		json.NewEncoder(w).Encode(500)
		Logus.Logus.Error(err)
		return
	}

	defer rows.Close()

	var notifications []NOTIFICATION_MODEL

	for rows.Next() {
		var notification NOTIFICATION_MODEL
		var createdAt []uint8

		err := rows.Scan(&notification.ID, &notification.Read, &notification.Title, &notification.ImageURL, &createdAt)
		if err != nil {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
			return
		}

		notification.CreatedAt, err = time.Parse("2006-01-02 15:04:05.000", string(createdAt))
		if err != nil {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
			return
		}

		notifications = append(notifications, notification)
	}

	if len(notifications) == 0 {
		notifications = []NOTIFICATION_MODEL{}
	}

	json.NewEncoder(w).Encode(notifications)
}
