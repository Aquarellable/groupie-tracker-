package handlers

import (
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	println(w, "Bienvenue sur mon API GroupieTrackers !")
}
