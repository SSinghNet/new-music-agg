package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SSinghNet/new-music-agg/internal/api"
	"github.com/SSinghNet/new-music-agg/internal/api/handler"
	"github.com/SSinghNet/new-music-agg/internal/service"
	"github.com/SSinghNet/new-music-agg/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Load dotenv failed")
	}

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

	r := api.NewRouter(releaseHandler)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
