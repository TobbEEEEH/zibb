package main

import (
	"log"
	"net/http"
	"time"

	"server/handlers"
	"server/internal/auth"
)

const port = ":42069"

func main() {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Load all templates into memory
	handlers.LoadTemplates()
	// Register all endpoints to mux
	handlers.RegisterHandlers(mux)

	auth.StartSessionCleanup(1 * time.Minute)

	log.Println("Listening on http://localhost" + port)
	log.Fatalln(http.ListenAndServe(port, mux))
}
