package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SSinghNet/new-music-agg/backend/internal/api"
	"github.com/SSinghNet/new-music-agg/backend/internal/api/handler"
	"github.com/SSinghNet/new-music-agg/backend/internal/service"
	"github.com/SSinghNet/new-music-agg/backend/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	st, err := store.NewPostgresStore(dsn)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer st.Close()

	releaseSvc := service.NewReleaseService(st)
	releaseHandler := handler.NewReleaseHandler(releaseSvc)

	artistSvc := service.NewArtistService(st)
	artistHandler := handler.NewArtistHandler(artistSvc)

	r := api.NewRouter(releaseHandler, artistHandler)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
