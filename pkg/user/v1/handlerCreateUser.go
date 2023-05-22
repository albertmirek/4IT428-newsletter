package v1

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"vse.com/4IT428/2023/newsletter/pkg/user/models"
)

func (h *Handlers) createUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		zap.Error(err)
		return
	}

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// Validate the password
	if req.Password == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	} else if len(req.Password) < 8 || len(req.Password) > 64 {
		http.Error(w, "Password must be between 8 and 64 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	uuID := uuid.New()

	newUser := models.User{
		ID:       uuID,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	_, err = h.UserRepository.CreateUser(r.Context(), newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		zap.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
