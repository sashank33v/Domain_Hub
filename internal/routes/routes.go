package routes

import (
	_ "domainhub/docs"
	"domainhub/internal/handlers"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(r chi.Router, domainHandler *handlers.DomainHandler) {
	r.Get("/health", handlers.HealthHandler)

	r.Post("/domains", domainHandler.Create)
	r.Get("/domains", domainHandler.GetAll)
	r.Get("/domains/{id}", domainHandler.GetByID)
	r.Put("/domains/{id}", domainHandler.Update)
	r.Delete("/domains/{id}", domainHandler.Delete)

	r.Get("/swagger/*", httpSwagger.Handler())
}
