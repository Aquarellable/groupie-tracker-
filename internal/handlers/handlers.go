package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"groupie-deezer-back/internal/api"
)

type Handlers struct {
	Deezer *api.DeezerClient
}

func New(c *api.DeezerClient) *Handlers {
	return &Handlers{Deezer: c}
}

/* ---------- MIDDLEWARE ---------- */

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				serverError(w, "panic recovered")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

/* ---------- HTML (templates) ---------- */

var tpl = template.Must(template.ParseGlob("templates/*.html"))

func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_ = tpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title": "Accueil",
		"View":  "home",
	})
}

func (h *Handlers) SearchPage(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, _, err := h.Deezer.SearchArtists(q, 10, 0)
	if err != nil {
		http.Error(w, "Erreur Deezer", 500)
		return
	}

	_ = tpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title":   "Résultats",
		"View":    "results",
		"Query":   q,
		"Artists": res.Data,
	})
}

/* ---------- JSON API ---------- */

func (h *Handlers) Docs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"routes": []any{
			map[string]string{"GET": "/api/docs", "desc": "routes + exemples"},
			map[string]string{"GET": "/api/search?type=artist&q=gims&limit=10&index=0", "desc": "search artist"},
			map[string]string{"GET": "/api/search?type=track&q=eminem&limit=10&index=0", "desc": "search track"},
			map[string]string{"GET": "/api/search?type=album&q=daft%20punk&limit=10&index=0", "desc": "search album"},
			map[string]string{"GET": "/api/artist/{id}", "desc": "artist details"},
			map[string]string{"GET": "/api/artist/{id}/top?limit=10", "desc": "artist top tracks"},
			map[string]string{"GET": "/api/artist/{id}/albums?limit=10&index=0", "desc": "artist albums"},
			map[string]string{"GET": "/api/album/{id}", "desc": "album details"},
			map[string]string{"GET": "/api/track/{id}", "desc": "track details"},
		},
	})
}

func getIntQuery(r *http.Request, key string, def int) int {
	v := strings.TrimSpace(r.URL.Query().Get(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func (h *Handlers) SearchAPI(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		badRequest(w, "param q manquant")
		return
	}

	typ := strings.TrimSpace(r.URL.Query().Get("type"))
	if typ == "" {
		typ = "artist"
	}

	limit := getIntQuery(r, "limit", 10)
	index := getIntQuery(r, "index", 0)

	switch typ {
	case "artist":
		res, _, err := h.Deezer.SearchArtists(q, limit, index)
		if err != nil {
			serverError(w, "deezer search artist failed")
			return
		}
		writeJSON(w, 200, res)

	case "track":
		res, _, err := h.Deezer.SearchTracks(q, limit, index)
		if err != nil {
			serverError(w, "deezer search track failed")
			return
		}
		writeJSON(w, 200, res)

	case "album":
		res, _, err := h.Deezer.SearchAlbums(q, limit, index)
		if err != nil {
			serverError(w, "deezer search album failed")
			return
		}
		writeJSON(w, 200, res)

	default:
		badRequest(w, "type invalide: artist|track|album")
	}
}

func (h *Handlers) ArtistRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/artist/")
	path = strings.Trim(path, "/")
	if path == "" {
		notFound(w, "id manquant")
		return
	}

	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		badRequest(w, "id invalide")
		return
	}

	if len(parts) == 1 {
		artist, _, err := h.Deezer.GetArtist(id)
		if err != nil {
			serverError(w, "deezer get artist failed")
			return
		}
		writeJSON(w, 200, artist)
		return
	}

	switch parts[1] {
	case "top":
		limit := getIntQuery(r, "limit", 10)
		top, _, err := h.Deezer.GetArtistTop(id, limit)
		if err != nil {
			serverError(w, "deezer artist top failed")
			return
		}
		writeJSON(w, 200, top)

	case "albums":
		limit := getIntQuery(r, "limit", 10)
		index := getIntQuery(r, "index", 0)
		albums, _, err := h.Deezer.GetArtistAlbums(id, limit, index)
		if err != nil {
			serverError(w, "deezer artist albums failed")
			return
		}
		writeJSON(w, 200, albums)

	default:
		notFound(w, "sous-route inconnue (top|albums)")
	}
}

func (h *Handlers) Album(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/album/")
	idStr = strings.Trim(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		badRequest(w, "id invalide")
		return
	}

	alb, _, err := h.Deezer.GetAlbum(id)
	if err != nil {
		serverError(w, "deezer get album failed")
		return
	}
	writeJSON(w, 200, alb)
}

func (h *Handlers) Track(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/track/")
	idStr = strings.Trim(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		badRequest(w, "id invalide")
		return
	}

	tr, _, err := h.Deezer.GetTrack(id)
	if err != nil {
		serverError(w, "deezer get track failed")
		return
	}
	writeJSON(w, 200, tr)
}