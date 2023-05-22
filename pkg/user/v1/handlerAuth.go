package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"vse.com/4IT428/2023/newsletter/pkg/user/models"
)

// Login handles the login request
func (h *Handlers) login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserRepository.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userEmail"] = user.Email
	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	jwtSecret := os.Getenv("JWT_SECRET")

	t, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	res := models.LoginResponse{
		Token: t,
	}

	json.NewEncoder(w).Encode(res)
}

//TODO implement refresh token
// RefreshTokenHandler is the handler for the refresh token endpoint.
// func (h *Handlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get the token from the Authorization header
// 	tokenString := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))

// 	// Parse and validate the token
// 	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("your_jwt_secret"), nil
// 	})

// 	if err != nil {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	// Check if the token is expired
// 	claims, ok := token.Claims.(*jwt.StandardClaims)
// 	if !ok || !token.Valid || time.Now().Unix() >= claims.ExpiresAt {
// 		http.Error(w, "Token is expired", http.StatusUnauthorized)
// 		return
// 	}

// 	// Generate a new token with a new expiration time
// 	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 		Subject:   claims.Subject,
// 		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
// 	})

// 	newTokenString, err := newToken.SignedString([]byte("your_jwt_secret"))
// 	if err != nil {
// 		http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Return the new token to the client
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"token": newTokenString,
// 	})
// }
