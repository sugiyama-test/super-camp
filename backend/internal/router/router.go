package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/handler"
	"github.com/sugiyamadaiki/super-camp/backend/internal/repository"
)

func New(pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	checklistRepo := repository.NewChecklistRepository(pool)
	checklistHandler := handler.NewChecklistHandler(checklistRepo)

	layoutRepo := repository.NewLayoutRepository(pool)
	layoutHandler := handler.NewLayoutHandler(layoutRepo)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck(pool))

		r.Route("/checklists", func(r chi.Router) {
			r.Get("/", checklistHandler.List)
			r.Post("/", checklistHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", checklistHandler.Get)
				r.Delete("/", checklistHandler.Delete)
				r.Post("/items", checklistHandler.AddItem)
				r.Put("/items/{itemID}", checklistHandler.UpdateItem)
				r.Delete("/items/{itemID}", checklistHandler.DeleteItem)
			})
		})

		r.Route("/layouts", func(r chi.Router) {
			r.Get("/", layoutHandler.List)
			r.Post("/", layoutHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", layoutHandler.Get)
				r.Put("/", layoutHandler.Update)
				r.Delete("/", layoutHandler.Delete)
			})
		})
	})

	return r
}
