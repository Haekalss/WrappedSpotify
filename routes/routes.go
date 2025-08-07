package routes

import (
	"net/http"

	"github.com/Haekalss/WrappedSpotify/handler"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", handler.SpotifyLogin)
	mux.HandleFunc("/auth/callback", handler.SpotifyCallback)
	mux.HandleFunc("/spotify/top-tracks", handler.GetTopTracks)
	mux.HandleFunc("/auth/refresh", handler.SpotifyRefreshToken)
	mux.HandleFunc("/spotify/top-artists", handler.GetTopArtists)
	mux.HandleFunc("/spotify/top-genres", handler.GetTopGenres)
	mux.HandleFunc("/spotify/wrapped", handler.GetWrappedData)

}
