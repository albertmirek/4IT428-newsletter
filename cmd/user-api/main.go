package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"vse.com/4IT428/2023/newsletter/pkg/user/repository"
	v1 "vse.com/4IT428/2023/newsletter/pkg/user/v1"
	v2 "vse.com/4IT428/2023/newsletter/pkg/user/v2"
	"vse.com/4IT428/2023/newsletter/shared/config"
	postgresql "vse.com/4IT428/2023/newsletter/shared/db/posgtresql"
)

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

	// Initialize the user repository
	userRepository := &repository.SQLUserRepository{
		DB: db,
	}

	// Create the router and register the handlers
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		zap.L().Log(logger.Level(), r.RequestURI)
		w.Write([]byte(`{"message": "pong (user-api)"}`))

	})

	v1.RegisterHandlers(r, userRepository)
	v2.RegisterHandlers(r, userRepository)

	// Start the HTTP server
	port := "8080"
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
		zap.Error(err)
	}
}
