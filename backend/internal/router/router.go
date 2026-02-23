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

	fireLogRepo := repository.NewFireLogRepository(pool)
	fireLogHandler := handler.NewFireLogHandler(fireLogRepo)

	mealPlanRepo := repository.NewMealPlanRepository(pool)
	mealPlanHandler := handler.NewMealPlanHandler(mealPlanRepo)

	gearRepo := repository.NewGearRepository(pool)
	gearHandler := handler.NewGearHandler(gearRepo)

	campsiteRepo := repository.NewCampsiteRepository(pool)
	campsiteHandler := handler.NewCampsiteHandler(campsiteRepo)

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

		r.Route("/fire-logs", func(r chi.Router) {
			r.Get("/", fireLogHandler.List)
			r.Post("/", fireLogHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", fireLogHandler.Get)
				r.Put("/", fireLogHandler.Update)
				r.Delete("/", fireLogHandler.Delete)
			})
		})

		r.Route("/meal-plans", func(r chi.Router) {
			r.Get("/", mealPlanHandler.List)
			r.Post("/", mealPlanHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", mealPlanHandler.Get)
				r.Put("/", mealPlanHandler.Update)
				r.Delete("/", mealPlanHandler.Delete)
			})
		})

		r.Route("/gears", func(r chi.Router) {
			r.Get("/", gearHandler.List)
			r.Post("/", gearHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", gearHandler.Get)
				r.Put("/", gearHandler.Update)
				r.Delete("/", gearHandler.Delete)
			})
		})

		r.Route("/campsites", func(r chi.Router) {
			r.Get("/", campsiteHandler.List)
			r.Post("/", campsiteHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", campsiteHandler.Get)
				r.Put("/", campsiteHandler.Update)
				r.Delete("/", campsiteHandler.Delete)
			})
		})
	})

	return r
}
