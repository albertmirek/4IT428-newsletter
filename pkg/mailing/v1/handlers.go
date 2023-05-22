package v1

import (
	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"vse.com/4IT428/2023/newsletter/pkg/mailing/repository"
	"vse.com/4IT428/2023/newsletter/shared/middleware"
	"vse.com/4IT428/2023/newsletter/shared/sendgrid"
)

type Handlers struct {
	MailingRepository repository.MailingRepository
	FirestoreClient   *firestore.Client
}

// RegisterHandlers registers the handlers for the user API
func RegisterHandlers(router *chi.Mux, mailingRepository repository.MailingRepository, firestoreClient *firestore.Client) {
	handlers := &Handlers{
		MailingRepository: mailingRepository,
		FirestoreClient:   firestoreClient,
	}

	router.With(middleware.AuthMiddleware).Post("/v1/{newsletterId}/send/{postId}", handlers.sendNewsletter)
}

func (h *Handlers) sendNewsletter(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.ContextUserKey).(string)
	newsletterID := chi.URLParam(r, "newsletterId")
	postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newsletterPost, err := h.MailingRepository.GetNewsletterWithPost(r.Context(), newsletterID, postID)
	if err != nil {
		log.Fatal(err)
		//log error
		zap.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if newsletterPost.AdminID.String() != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get all users subscribed to newsletter
	users, err := h.FirestoreClient.Collection("newsletters").Doc(newsletterID).Collection("subscribers").Documents(r.Context()).GetAll()
	if err != nil {
		log.Fatal(err)
		zap.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send newsletter to all users
	for _, user := range users {
		userData := user.Data()
		userEmail := userData["email"].(string)
		userToken := userData["token"].(string)

		newsletterEmail := sendgrid.NewsletterEmail{
			To:               userEmail,
			UnsubscribeToken: userToken,
			Subject:          newsletterPost.NameOfNewsletter,
			Body:             newsletterPost.Body,
		}

		// send email
		_, err := sendgrid.SendNewsletterEmail(newsletterEmail)
		if err != nil {
			log.Fatal("Senrgrid", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	log.Println("Emails sent")

	w.WriteHeader(http.StatusOK)

}
