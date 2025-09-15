// Package handlers serves API
package handlers

import (
	"log"
	"net/http"

	"server/internal/auth"
	"server/models"
	"server/storage"
)

// LoginHandler serves /login path
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	storage.UpdateRequestCounter("login")
	errorMessage := []byte(`<p class="error-message">Wrong username or password</p>`)

	switch r.Method {
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, ok := storage.Users.Load(username)
		if !ok {
			w.Write([]byte(errorMessage))
			return
		}

		if u, ok := user.(models.User); ok {
			valid, _ := auth.CheckPassword(u.Hash, password)
			if !valid {
				w.Write(errorMessage)
				return
			}

			log.Printf("User %s logged in", username)
			auth.CreateSession(w, username)

			w.Header().Set("HX-Redirect", "/profile")
			w.WriteHeader(http.StatusOK)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
