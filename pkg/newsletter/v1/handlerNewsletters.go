package v1

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"vse.com/4IT428/2023/newsletter/pkg/newsletter/models"
	"vse.com/4IT428/2023/newsletter/shared/middleware"
	"vse.com/4IT428/2023/newsletter/shared/sendgrid"
)

func (h *Handlers) getNewsletters(w http.ResponseWriter, r *http.Request) {

	newsletters, err := h.NewsletterRepository.GetNewsletters(r.Context())
	if err != nil {
		http.Error(w, "Failed to get newsletters", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsletters)
}

func (h *Handlers) createNewsletter(w http.ResponseWriter, r *http.Request) {

	var req models.CreateNewsletterRequest

	userUUID, err := uuid.Parse(r.Context().Value(middleware.ContextUserKey).(string))
	if err != nil {
		http.Error(w, "Failed to parse user UUID", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		zap.Error(err)
		return
	}
	uuID := uuid.New()

	newNewsletter := models.Newsletter{
		ID:     uuID,
		UserId: userUUID,
		Name:   req.Name,
	}

	createNewsletter, err := h.NewsletterRepository.CreateNewsletter(r.Context(), newNewsletter)
	if err != nil {
		http.Error(w, "Failed to create newsletter", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createNewsletter)
}

func (h *Handlers) creatNewsletterPost(w http.ResponseWriter, r *http.Request) {

	var req models.CreateNewsletterPost

	userID := r.Context().Value(middleware.ContextUserKey).(string)
	newsletterUUID, err := uuid.Parse(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		zap.Error(err)
		return
	}

	newsletter, err := h.NewsletterRepository.GetNewsletterById(r.Context(), newsletterUUID.String())
	if err != nil {
		http.Error(w, "Newsletter not found", http.StatusNotFound)
		zap.Error(err)
		return
	}

	if newsletter.UserId.String() != userID {
		http.Error(w, "You are not admin of this newsletter", http.StatusForbidden)
		zap.Error(err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		zap.Error(err)
		return
	}

	newNewsletterPost := models.Post{
		NewsletterID: newsletterUUID,
		Heading:      req.Heading,
		Body:         req.Body,
	}

	_, err = h.NewsletterRepository.CreateNewsletterPost(r.Context(), newNewsletterPost)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) updateNewsletter(w http.ResponseWriter, r *http.Request) {

	var req models.UpdateNewsletterRequest

	userUUID, err := uuid.Parse(r.Context().Value(middleware.ContextUserKey).(string))
	if err != nil {
		http.Error(w, "Failed to parse user UUID", http.StatusInternalServerError)
		zap.Error(err)
		return
	}
	newsletterUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		zap.Error(err)
		return
	}

	newsletter, err := h.NewsletterRepository.GetNewsletterById(r.Context(), newsletterUUID.String())
	if err != nil {
		http.Error(w, "Newsletter not found", http.StatusNotFound)
		zap.Error(err)
		return
	}

	if newsletter.UserId != userUUID {
		http.Error(w, "You are not admin of this newsletter", http.StatusForbidden)
		zap.Error(err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		zap.Error(err)
		return
	}

	newNewsletter := models.Newsletter{
		ID:     newsletterUUID,
		UserId: userUUID,
		Name:   req.Name,
	}

	_, err = h.NewsletterRepository.UpdateNewsletterById(r.Context(), newNewsletter)
	if err != nil {
		http.Error(w, "Failed to update newsletter", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) deleteNewsletter(w http.ResponseWriter, r *http.Request) {

	userUUID, err := uuid.Parse(r.Context().Value(middleware.ContextUserKey).(string))
	if err != nil {
		http.Error(w, "Failed to parse user UUID", http.StatusInternalServerError)
		zap.Error(err)
		return
	}
	newsletterUUID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid newsletter ID", http.StatusBadRequest)
		zap.Error(err)
		return
	}

	newsletter, err := h.NewsletterRepository.GetNewsletterById(r.Context(), newsletterUUID.String())
	if err != nil {
		http.Error(w, "Newsletter not found", http.StatusNotFound)
		zap.Error(err)
		return
	}

	if newsletter.UserId != userUUID {
		http.Error(w, "You are not admin of this newsletter", http.StatusForbidden)
		zap.Error(err)
		return
	}

	firestoreClient := h.FirestoreClient
	subscribersRef := firestoreClient.Collection("newsletters").Doc(newsletterUUID.String()).Collection("subscribers")

	users, err := subscribersRef.Documents(r.Context()).GetAll()

	if err != nil {
		http.Error(w, "Failed to get subscribers", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	err = h.NewsletterRepository.DeletePostsByNewsletterId(r.Context(), newsletterUUID.String())
	if err != nil {
		http.Error(w, "Failed to delete posts", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	err = h.NewsletterRepository.DeleteNewsletterById(r.Context(), newsletterUUID.String())
	if err != nil {
		http.Error(w, "Failed to delete newsletter db", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	if err := deleteCollection(r.Context(), firestoreClient, subscribersRef, 50); err != nil {
		log.Fatalf("Failed to delete 'subscribers' subcollection: %v", err)
	}
	if _, err := firestoreClient.Collection("newsletters").Doc(newsletterUUID.String()).Delete(r.Context()); err != nil {
		log.Fatalf("Failed to delete newsletter document: %v", err)
	}

	err = h.NewsletterRepository.DeleteNewsletterById(r.Context(), newsletterUUID.String())
	if err != nil {
		http.Error(w, "Failed to delete newsletter", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	for _, user := range users {
		userData := user.Data()
		userEmail := userData["email"].(string)

		newsletterEmail := sendgrid.ConfirmationEmail{
			To:             userEmail,
			NewsletterName: newsletter.Name,
		}

		_, err = sendgrid.SendNewsletterDeletionEmail(newsletterEmail)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			zap.Error(err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCollection(ctx context.Context, client *firestore.Client, collRef *firestore.CollectionRef, batchSize int) error {
	for {
		// Retrieve a batch of documents
		iter := collRef.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Begin a new batch.
		batch := client.Batch()

		// Iterate through the documents, adding
		// a delete operation for each one to the batch.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete, the process is done.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
