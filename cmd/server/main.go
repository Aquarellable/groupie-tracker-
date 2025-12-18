package main

import (
	"log"
	"net/http"
	"time"

	"groupie-deezer-back/internal/api"
	"groupie-deezer-back/internal/handlers"
)

func main() {
	deezerClient := api.NewDeezerClient("https://api.deezer.com", 10*time.Second)
	h := handlers.New(deezerClient)

	mux := http.NewServeMux()

	// HTML
	mux.HandleFunc("/", h.HomePage)
	mux.HandleFunc("/search", h.SearchPage)

	// JSON API
	mux.HandleFunc("/api/docs", h.Docs)
	mux.HandleFunc("/api/search", h.SearchAPI)
	mux.HandleFunc("/api/artist/", h.ArtistRouter)
	mux.HandleFunc("/api/album/", h.Album)
	mux.HandleFunc("/api/track/", h.Track)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlers.RecoverMiddleware(mux),
	}

	log.Println("✅ Server lancé sur http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}