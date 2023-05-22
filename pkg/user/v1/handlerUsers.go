package v1

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"vse.com/4IT428/2023/newsletter/pkg/user/models"
	"vse.com/4IT428/2023/newsletter/shared/middleware"

	"go.uber.org/zap"
)

func (h *Handlers) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserRepository.GetAllUsers(r.Context())
	if err != nil {
		zap.Error(err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}

func (h *Handlers) updateUserPassword(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateUserPasswordRequest

	userEmail := r.Context().Value(middleware.ContextUserEmailKey).(string)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.Error(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	} else if len(req.NewPassword) < 8 || len(req.NewPassword) > 64 {
		http.Error(w, "Password must be between 8 and 64 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	err = h.UserRepository.UpdateUserPassword(r.Context(), string(hashedPassword), userEmail)

	w.WriteHeader(http.StatusOK)
}
