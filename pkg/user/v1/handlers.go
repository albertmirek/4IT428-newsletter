package v1

import (
	"net/http"
	"vse.com/4IT428/2023/newsletter/shared/middleware"

	"github.com/go-chi/chi"
	"vse.com/4IT428/2023/newsletter/pkg/user/repository"
)

// Handlers contains all handlers for the user API
type Handlers struct {
	UserRepository repository.UserRepository
}

// RegisterHandlers registers the handlers for the user API
func RegisterHandlers(router *chi.Mux, userRepository repository.UserRepository) {
	handlers := &Handlers{
		UserRepository: userRepository,
	}

	router.Get("/v1/ping", pingHandler)

	router.Post("/v1/login", handlers.login)

	//router.Get("/v1/users", handlers.getAllUsers)

	router.Post("/v1", handlers.createUser)

	router.With(middleware.AuthMiddleware).Patch("/v1", handlers.updateUserPassword)

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "healthy (v1)"}`))
}
