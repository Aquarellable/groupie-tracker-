package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type DeezerClient struct {
	BaseURL string
	Client  *http.Client
}

func NewDeezerClient(baseURL string, timeout time.Duration) *DeezerClient {
	return &DeezerClient{
		BaseURL: baseURL,
		Client: &http.Client{Timeout: timeout},
	}
}

/* ---------- MODELS (subset utile) ---------- */

type Artist struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	PictureMedium string `json:"picture_medium"`
	NbAlbum       int    `json:"nb_album"`
	NbFan         int    `json:"nb_fan"`
	Tracklist     string `json:"tracklist"`
}

type Album struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Cover       string `json:"cover"`
	CoverMedium string `json:"cover_medium"`
	ReleaseDate string `json:"release_date"`
	Tracklist   string `json:"tracklist"`
}

type Track struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Duration int    `json:"duration"`
	Rank     int    `json:"rank"`
	Preview  string `json:"preview"`
	Artist   Artist `json:"artist"`
	Album    Album  `json:"album"`
}

type ListResponse[T any] struct {
	Data  []T    `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next,omitempty"`
}

/* ---------- INTERNAL ---------- */

func (c *DeezerClient) get(path string, target any) (int, error) {
	resp, err := c.Client.Get(c.BaseURL + path)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("deezer status=%d", resp.StatusCode)
	}

	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(target)
}

/* ---------- SEARCH ---------- */

func (c *DeezerClient) SearchArtists(q string, limit, index int) (*ListResponse[Artist], int, error) {
	if limit <= 0 {
		limit = 10
	}
	if index < 0 {
		index = 0
	}
	qs := url.QueryEscape(q)
	path := fmt.Sprintf("/search/artist?q=%s&limit=%d&index=%d", qs, limit, index)

	var res ListResponse[Artist]
	code, err := c.get(path, &res)
	return &res, code, err
}

func (c *DeezerClient) SearchTracks(q string, limit, index int) (*ListResponse[Track], int, error) {
	if limit <= 0 {
		limit = 10
	}
	if index < 0 {
		index = 0
	}
	qs := url.QueryEscape(q)
	path := fmt.Sprintf("/search/track?q=%s&limit=%d&index=%d", qs, limit, index)

	var res ListResponse[Track]
	code, err := c.get(path, &res)
	return &res, code, err
}

func (c *DeezerClient) SearchAlbums(q string, limit, index int) (*ListResponse[Album], int, error) {
	if limit <= 0 {
		limit = 10
	}
	if index < 0 {
		index = 0
	}
	qs := url.QueryEscape(q)
	path := fmt.Sprintf("/search/album?q=%s&limit=%d&index=%d", qs, limit, index)

	var res ListResponse[Album]
	code, err := c.get(path, &res)
	return &res, code, err
}

/* ---------- DETAILS ---------- */

func (c *DeezerClient) GetArtist(id int) (*Artist, int, error) {
	var a Artist
	code, err := c.get(fmt.Sprintf("/artist/%d", id), &a)
	return &a, code, err
}

func (c *DeezerClient) GetArtistTop(id int, limit int) (*ListResponse[Track], int, error) {
	if limit <= 0 {
		limit = 10
	}
	var res ListResponse[Track]
	code, err := c.get(fmt.Sprintf("/artist/%d/top?limit=%d", id, limit), &res)
	return &res, code, err
}

func (c *DeezerClient) GetArtistAlbums(id int, limit, index int) (*ListResponse[Album], int, error) {
	if limit <= 0 {
		limit = 10
	}
	if index < 0 {
		index = 0
	}
	var res ListResponse[Album]
	code, err := c.get(fmt.Sprintf("/artist/%d/albums?limit=%d&index=%d", id, limit, index), &res)
	return &res, code, err
}

func (c *DeezerClient) GetAlbum(id int) (*Album, int, error) {
	var a Album
	code, err := c.get(fmt.Sprintf("/album/%d", id), &a)
	return &a, code, err
}

func (c *DeezerClient) GetTrack(id int) (*Track, int, error) {
	var t Track
	code, err := c.get(fmt.Sprintf("/track/%d", id), &t)
	return &t, code, err
}