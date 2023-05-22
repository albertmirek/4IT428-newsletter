package v2

import (
	"net/http"

	"github.com/go-chi/chi"
	"vse.com/4IT428/2023/newsletter/pkg/user/repository"
)

func RegisterHandlers(router *chi.Mux, userRepository repository.UserRepository) {
	router.Get("/v2/ping", pingHandler)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "pong (v2)"}`))
}
