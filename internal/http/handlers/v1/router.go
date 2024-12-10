package handlers

import (
	"effective-mobile-test/internal/http/middlewares/pagination"
	"effective-mobile-test/internal/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewRouter(log *slog.Logger, r *chi.Mux, sluc *usecases.SongLibrary) {
	r.Use(
		middleware.RequestID,
		middleware.Recoverer,
		middleware.URLFormat,
	)

	sl := newSongLibrary(sluc, log)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/songs", func(r chi.Router) {
			r.
				With(pagination.SetPaginationContextMiddleware).
				Get("/", sl.getList)

			r.Post("/", sl.create)
			r.Patch("/", sl.update)
			r.Delete("/", sl.delete)

			r.Route("/text", func(r chi.Router) {
				r.
					With(pagination.SetPaginationContextMiddleware).
					Get("/", sl.getText)
			})
		})
	})

	r.Route("/info", func(r chi.Router) {
		r.Get("/", sl.get)
	})
}
