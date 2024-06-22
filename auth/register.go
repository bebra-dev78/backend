package auth

import (
	Logus "backend/logging"
	"database/sql"
	"encoding/json"
	"net/http"
	"unicode/utf8"

	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	if name == "" || email == "" || password == "" {
		json.NewEncoder(w).Encode(400)
		return
	}

	if utf8.RuneCountInString(name) < 3 || utf8.RuneCountInString(name) > 13 || utf8.RuneCountInString(email) < 1 || utf8.RuneCountInString(email) > 255 {
		json.NewEncoder(w).Encode(nil)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	if _, err := db.Exec("INSERT INTO users (id, email, password, name) VALUES (?, ?, ?, ?)", cuid.New(), email, hash, name); err != nil {
		json.NewEncoder(w).Encode(500)
		Logus.Logus.Error(err)
		return
	}

	json.NewEncoder(w).Encode(200)
}
