package v1

import (
	"github.com/go-chi/chi"
	"net/http"
	"vse.com/4IT428/2023/newsletter/shared/middleware"
	"vse.com/4IT428/2023/newsletter/shared/sendgrid"
)

func (h *Handlers) subscribeToNewsletter(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserKey).(string)
	userEmail := r.Context().Value(middleware.ContextUserEmailKey).(string)
	newsletterID := chi.URLParam(r, "id")

	newsletter, err := h.NewsletterRepository.GetNewsletterById(r.Context(), newsletterID)
	if err != nil {
		http.Error(w, "Newsletter not found", http.StatusNotFound)
		return
	}

	userNewsletterToken, err := sendgrid.EncryptUserNewsletterToken(userID, userEmail, newsletterID)
	if err != nil {
		http.Error(w, "Failed generate pairing token", http.StatusInternalServerError)
		return
	}

	_, err = h.FirestoreClient.Collection("newsletters").Doc(newsletterID).Collection("subscribers").Doc(userID).Set(r.Context(), map[string]interface{}{
		"email": userEmail,
		"token": userNewsletterToken,
	})

	if err != nil {
		http.Error(w, "Failed to subscribe to newsletter", http.StatusInternalServerError)
		return
	}

	confirmationEmail := sendgrid.ConfirmationEmail{
		To:             userEmail,
		NewsletterName: newsletter.Name,
	}

	_, err = sendgrid.SendSubscribeConfirmationEmail(confirmationEmail, userNewsletterToken)
	if err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *Handlers) unsubscribeFromNewsletter(w http.ResponseWriter, r *http.Request) {
	encodedToken := chi.URLParam(r, "token")
	data, err := sendgrid.GetValuesFromEncryptedToken(encodedToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	_, err = h.FirestoreClient.Collection("newsletters").Doc(data.NewsletterID).Collection("subscribers").Doc(data.UserID).Delete(r.Context())
	if err != nil {
		http.Error(w, "Failed to unsubscribe from newsletter", http.StatusInternalServerError)
		return
	}

	newsletter, err := h.NewsletterRepository.GetNewsletterById(r.Context(), data.NewsletterID)
	if err != nil {
		http.Error(w, "Newsletter not found", http.StatusNotFound)
		return
	}

	confirmationEmail := sendgrid.ConfirmationEmail{
		To:             data.UserEmail,
		NewsletterName: newsletter.Name,
	}

	_, err = sendgrid.SendUnsubscribeConfirmationEmail(confirmationEmail)
	if err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
