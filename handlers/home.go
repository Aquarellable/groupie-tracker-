package handlers

import (
	"net/http"
	"text/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/templates/index.html"))
	tmpl.Execute(w, nil)
}
