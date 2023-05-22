package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"vse.com/4IT428/2023/newsletter/pkg/mailing/repository"
	v1 "vse.com/4IT428/2023/newsletter/pkg/mailing/v1"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"vse.com/4IT428/2023/newsletter/shared/config"
	firestoreDb "vse.com/4IT428/2023/newsletter/shared/db/firestore"
	postgresql "vse.com/4IT428/2023/newsletter/shared/db/posgtresql"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.Error(err)
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the PostgreSQL database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgresql.NewConnectionPool(ctx, cfg.Database)
	if err != nil {
		zap.Error(err)
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	mailingRepository := &repository.SQLMailingRepository{
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
		w.Write([]byte(`{"message": "pong (mailing-api)"}`))
	})

	/*r.Get("/testmail", func(w http.ResponseWriter, r *http.Request) {
		from := mail.NewEmail("Example User", "am7642939@gmail.com")
		subject := "Sending with SendGrid is Fun"
		to := mail.NewEmail("Example User", "test@example.com")
		plainTextContent := "and easy to do anywhere, even with Go"
		htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

		log.Println("KEY", os.Getenv("SENDGRID_API_KEY"))
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	})*/

	v1.RegisterHandlers(r, mailingRepository, firestoreClient)
	// v2.RegisterHandlers(r, mailingRepository, firestoreClient)

	// Start the HTTP server
	port := "8080"
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil && err != http.ErrServerClosed {
		zap.Error(err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
