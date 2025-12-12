package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/handlers"
	"net/http"
	"time"
)

// Structure qui correspond au JSON retourné par l'API Groupie Trackers
type Artist struct {
	ID         int      `json:"id"`
	Image      string   `json:"image"`
	Name       string   `json:"name"`
	Members    []string `json:"members"`
	Creation   int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
}

// Handler : récupère les données depuis l’API externe et les renvoie
func getArtists(w http.ResponseWriter, r *http.Request) {

	// URL de l'API GroupieTrackers
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Requête GET
	resp, err := client.Get(apiURL)
	if err != nil {
		http.Error(w, "Erreur lors de la requête externe", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Décode le JSON dans un slice de Artist
	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		http.Error(w, "Erreur : JSON impossible à lire", http.StatusInternalServerError)
		return
	}

	// Répond en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// Fonction principale
func main() {
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/artists", getArtists)

	fmt.Println("API ready !")
	fmt.Println("Listening at : http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

type Artiste struct {
	Nom string `json:"nom"`
}
