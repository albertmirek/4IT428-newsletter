package v1

import (
	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"vse.com/4IT428/2023/newsletter/pkg/newsletter/repository"
	"vse.com/4IT428/2023/newsletter/shared/middleware"
)

type Handlers struct {
	NewsletterRepository repository.NewsletterRepository
	FirestoreClient      *firestore.Client
}

func RegisterHandlers(router *chi.Mux, newsletterRepository repository.NewsletterRepository, firestoreClient *firestore.Client) {
	handlers := &Handlers{
		NewsletterRepository: newsletterRepository,
		FirestoreClient:      firestoreClient,
	}

	router.Get("/v1/newsletter", handlers.getNewsletters)
	router.With(middleware.AuthMiddleware).Post("/v1/newsletter", handlers.createNewsletter)

	router.With(middleware.AuthMiddleware).Post("/v1/newsletter/{id}/post", handlers.creatNewsletterPost)
	router.With(middleware.AuthMiddleware).Put("/v1/newsletter/{id}", handlers.updateNewsletter)
	router.With(middleware.AuthMiddleware).Delete("/v1/newsletter/{id}", handlers.deleteNewsletter)

	router.With(middleware.AuthMiddleware).Post("/v1/newsletter/{id}/subscribe", handlers.subscribeToNewsletter)
	router.Get("/v1/newsletter/{token}/unsubscribe", handlers.unsubscribeFromNewsletter)

}
