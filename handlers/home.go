package handlers

import (
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	Fprintln(w, "Bienvenue sur mon API GroupieTrackers !")
}
