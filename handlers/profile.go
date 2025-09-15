package handlers

import (
	"fmt"
	"net/http"

	"server/internal/auth"
	"server/models"
	"server/storage"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	storage.UpdateRequestCounter("profile")
	tmpl, err := Templates.Load("profile.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Println("Could not load template:", err)
		return
	}
	if username, ok := auth.UserIDFromContext(r.Context()); ok {
		if err := tmpl.Execute(w, models.User{
			Username: username,
		}); err != nil {
			http.Error(w, "Could not load page", http.StatusInternalServerError)
			fmt.Println("Could not load template:", err)
		}
		return
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
