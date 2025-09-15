package handlers

import (
	"log"
	"net/http"

	"server/internal/auth"
	"server/models"
	"server/storage"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	storage.UpdateRequestCounter("signup")
	signup, err := Templates.Load("signupForm.html")
	if err != nil {
		http.Error(w, "Issue loading HTML", http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := signup.Execute(w, models.Message{Message: ""}); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
			return
		}
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		if username == "" {
			signup.Execute(w, models.Message{
				Message: "Invalid username",
			})
			return
		}
		if password == "" {
			signup.Execute(w, models.Message{
				Message: "Invalid password",
			})
			return
		}
		if email == "" {
			signup.Execute(w, models.Message{
				Message: "Invalid email",
			})
			return
		}

		if _, exists := storage.Users.Load(username); exists {
			signup.Execute(w, models.Message{
				Message: "Username is taken",
			})
			return
		}

		hash, _ := auth.HashPassword(password)
		storage.Users.Store(username, models.User{Username: username, Email: email, Hash: hash})
		auth.CreateSession(w, username)
		w.Header().Set("HX-Redirect", "/profile")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}
