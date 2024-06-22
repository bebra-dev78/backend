package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	Logus "backend/logging"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	secret := r.Header.Get("Secret")
	if secret != "bebra" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	if email == "" || password == "" {
		json.NewEncoder(w).Encode(400)
		return
	}

	var userPassword string

	if err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&userPassword); err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(404)
		} else {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}

	json.NewEncoder(w).Encode(200)
}
