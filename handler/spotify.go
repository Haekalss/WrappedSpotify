package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Haekalss/WrappedSpotify/utils"
)

func GetTopTracks(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		http.Error(w, "access_token is required", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?limit=10&time_range=long_term", nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get top tracks", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(body), http.StatusInternalServerError)
		return
	}
	var result struct {
		Items []interface{} `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse tracks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Items)
}

func GetTopArtists(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		http.Error(w, "access_token is required", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/artists?limit=10&time_range=long_term", nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get top artists", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(body), http.StatusInternalServerError)
		return
	}
	var result struct {
		Items []interface{} `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse artists", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Items)
}

func GetTopGenres(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		http.Error(w, "access_token is required", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/artists?limit=20&time_range=long_term", nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get artists for genre analysis", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(body), http.StatusInternalServerError)
		return
	}

	var result struct {
		Items []struct {
			Genres []string `json:"genres"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse genre data", http.StatusInternalServerError)
		return
	}

	genreCount := make(map[string]int)
	for _, item := range result.Items {
		for _, genre := range item.Genres {
			genreCount[genre]++
		}
	}

	// Convert to array of { genre, count }
	type GenreData struct {
		Genre string `json:"genre"`
		Count int    `json:"count"`
	}
	var genres []GenreData
	for genre, count := range genreCount {
		genres = append(genres, GenreData{Genre: genre, Count: count})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func GetWrappedData(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	refreshToken := r.URL.Query().Get("refresh_token")
	if accessToken == "" || refreshToken == "" {
		http.Error(w, "access_token and refresh_token are required", http.StatusBadRequest)
		return
	}

	// Fetch using auto refresh logic
	topArtists, newToken, err := utils.CallSpotifyWithRefresh("https://api.spotify.com/v1/me/top/artists?limit=10&time_range=long_term", accessToken, refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topTracks, _, err := utils.CallSpotifyWithRefresh("https://api.spotify.com/v1/me/top/tracks?limit=10&time_range=long_term", newToken, refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Genre parsing
	var artistResult struct {
		Items []struct {
			Genres []string `json:"genres"`
		} `json:"items"`
	}
	_ = json.Unmarshal(topArtists, &artistResult)
	genreCount := map[string]int{}
	for _, a := range artistResult.Items {
		for _, g := range a.Genres {
			genreCount[g]++
		}
	}
	type GenreStat struct {
		Genre string `json:"genre"`
		Count int    `json:"count"`
	}
	var genres []GenreStat
	for genre, count := range genreCount {
		genres = append(genres, GenreStat{Genre: genre, Count: count})
	}

	result := map[string]interface{}{
		"top_artists":  json.RawMessage(topArtists),
		"top_tracks":   json.RawMessage(topTracks),
		"top_genres":   genres,
		"access_token": newToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
