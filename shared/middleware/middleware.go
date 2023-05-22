package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID    uuid.UUID
	UserEmail string
	jwt.StandardClaims
}

type ContextKey string

const ContextUserKey ContextKey = "userID"
const ContextUserEmailKey ContextKey = "userEmail"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT string from the header
		tokenStr := r.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Validate the token and extract the claims
		claims := &Claims{}

		jwtSecret := os.Getenv("JWT_SECRET")

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set userID in context for further use in handlers
		ctx := context.WithValue(r.Context(), ContextUserKey, claims.UserID.String())
		ctx = context.WithValue(ctx, ContextUserEmailKey, claims.UserEmail)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
