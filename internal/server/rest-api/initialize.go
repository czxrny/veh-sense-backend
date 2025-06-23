package restapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitializeAndStart() error {
	router := chi.NewRouter()
	initializeHandlers(router)
	return http.ListenAndServe(":8080", router)
}

func initializeHandlers(router *chi.Mux) {
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}
