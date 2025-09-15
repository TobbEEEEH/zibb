package handlers

import (
	"log"
	"net/http"

	"server/models"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	index, err := Templates.Load("index.html")
	if err != nil {
		http.Error(w, "Template was not found", http.StatusInternalServerError)
		log.Println("Template was not found:", err)
		return
	}

	if err := index.Execute(w, models.PageData{Theme: "dark", Title: "Rootsida"}); err != nil {
		http.Error(w, "Template was not found", http.StatusInternalServerError)
		log.Println("Error executing template: ", err)
		return
	}
}
