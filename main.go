package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type PageData struct {
	Theme string
}

func main() {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Parse templates once at startup
	index := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

	// Serve homepage with theme data
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{Theme: "dark"}
		err := index.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
	})

	// Dummy /submit handler for HTMX form submission
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// You can parse form data here:
		email := r.FormValue("email")

		// Respond with a simple confirmation snippet
		w.Header().Set("Content-Type", "text/html")
		response := `<p>Thanks ` + template.HTMLEscapeString(email) + `, but this wont do anything yet. <strong></strong></p>`
		w.Write([]byte(response))
	})

	port := ":42069"
	log.Println("Listening on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}
