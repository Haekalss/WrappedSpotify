package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Haekalss/WrappedSpotify/routes"
	"github.com/Haekalss/WrappedSpotify/utils"
)

func main() {
	utils.LoadEnv()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	handler := utils.CORSMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
