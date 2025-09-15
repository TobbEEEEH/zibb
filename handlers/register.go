package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"server/internal/auth"
	"server/models"
)

//go:embed templates/*
var templateFS embed.FS

var Templates models.Templates

func LoadTemplates() error {
	Templates = models.Templates{}
	loadTemplate("templates")
	return nil
}

func loadTemplate(path string) {
	entries, err := templateFS.ReadDir(path)
	if err != nil {
		log.Fatal("Was not able to read template directory")
	}

	for _, e := range entries {
		if e.IsDir() {
			newPath := filepath.Join(path, e.Name())
			loadTemplate(newPath)
		} else {
			if _, err := Templates.Load(e.Name()); err == nil {
				continue
			}
			tmpl := template.Must(template.ParseFS(templateFS, path+"/"+e.Name()))
			Templates.Add(e.Name(), tmpl)
			log.Printf("Added '%s' to memory from path %s\n", e.Name(), path+"/"+e.Name())
		}
	}
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/signup", SignupHandler)

	mux.Handle("/profile", auth.CookieAuthMiddleware(http.HandlerFunc(ProfileHandler)))
}
