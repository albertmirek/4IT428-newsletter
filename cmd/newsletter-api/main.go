package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"vse.com/4IT428/2023/newsletter/pkg/newsletter/repository"
	v1 "vse.com/4IT428/2023/newsletter/pkg/newsletter/v1"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"vse.com/4IT428/2023/newsletter/shared/config"
	firestoreDb "vse.com/4IT428/2023/newsletter/shared/db/firestore"
	postgresql "vse.com/4IT428/2023/newsletter/shared/db/posgtresql"
)

// TODO, inspiration from user-api
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		zap.Error(err)
	}

	// Connect to the PostgreSQL database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgresql.NewConnectionPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		zap.Error(err)
	}
	defer db.Close()

	newsletterRepository := &repository.SQLNewsletterRepository{
		DB: db,
	}

	firestoreClient, err := firestoreDb.InitializeFirebase()
	if err != nil {
		zap.Error(err)
		log.Fatalf("Failed to connect to the firestore: %v", err)
	}
	defer firestoreClient.Close()

	// Create the router and register the handlers
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "pong (newsletterrrrrrr-api)"}`))
	})

	v1.RegisterHandlers(r, newsletterRepository, firestoreClient)
	// v2.RegisterHandlers(r, userRepository)

	// Start the HTTP server
	port := "8080"
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
		zap.Error(err)
	}
}
