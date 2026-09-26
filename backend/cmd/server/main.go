package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sezzle-calculator/backend/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           newHandler(os.Getenv("STATIC_DIR")),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// newHandler serves the API and, when staticDir is set, the built frontend
// from the same origin.
func newHandler(staticDir string) http.Handler {
	apiHandler := api.NewHandler()
	if staticDir == "" {
		return apiHandler
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}
