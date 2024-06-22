package resourses

import (
	"database/sql"
	"encoding/json"
	"net/http"

	Logus "backend/logging"
)

type PANEL_MODEL struct {
	ID      uint    `db:"id"`
	Title   string  `db:"title"`
	Widgets *string `db:"widgets"`
}

type Response struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Widgets []int8 `json:"widgets"`
}

func Panels(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	rows, err := db.Query("SELECT id, title, widgets FROM panels WHERE uid = ?", uid)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	defer rows.Close()

	var panels []Response

	for rows.Next() {
		var panel PANEL_MODEL

		if err := rows.Scan(&panel.ID, &panel.Title, &panel.Widgets); err != nil {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
			return
		}

		var widgets []int8

		if panel.Widgets != nil {
			if err := json.Unmarshal([]byte(*panel.Widgets), &widgets); err != nil {
				json.NewEncoder(w).Encode(501)
				Logus.Logus.Error(err)
				return
			}
		}

		response := Response{
			ID:      panel.ID,
			Title:   panel.Title,
			Widgets: widgets,
		}
		panels = append(panels, response)
	}

	if len(panels) == 0 {
		panels = []Response{}
	}

	json.NewEncoder(w).Encode(panels)
}
