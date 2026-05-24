package main

import (
	"context"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/SSinghNet/new-music-agg/backend/internal/scraper"
	"github.com/SSinghNet/new-music-agg/backend/internal/store"

	"github.com/joho/godotenv"
)

const workers = 10
const pages = 50

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	st, err := store.NewPostgresStore(dsn)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer st.Close()

	ctx := context.Background()
	releases := scraper.Run(ctx, pages)

	var count, errCount atomic.Int64
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for r := range releases {
		r := r
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			if err := st.Upsert(ctx, r); err != nil {
				log.Printf("upsert %q: %v", r.Name, err)
				errCount.Add(1)
				return
			}
			count.Add(1)
		}()
	}

	wg.Wait()
	log.Printf("done: upserted=%d errors=%d", count.Load(), errCount.Load())
}
