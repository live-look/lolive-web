package server

import (
	"html/template"
	"log"
	"net/http"
)

// HomePage handles main page of site
func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("app").ParseFiles(
		"web/template/layout.html",
		"web/template/index.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "layout.html", nil)
}
