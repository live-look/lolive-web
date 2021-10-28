package server

import (
	"html/template"
	"log"
	"net/http"
)

func BroadcastNew(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("app").ParseFiles(
		"web/template/layout.html",
		"web/template/broadcasts/new.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "layout.html", nil)
}
