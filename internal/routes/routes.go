package routes

import (
	_ "domainhub/docs"
	"domainhub/internal/handlers"
	"domainhub/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(r chi.Router, domainHandler *handlers.DomainHandler, authhandler *handlers.AuthHandler, chathandler *handlers.ChatHandler, jwtSecret string) {
	r.Get("/health", handlers.HealthHandler)

	r.Post("/auth/register", authhandler.Register)
	r.Post("/auth/login", authhandler.Login)
	r.Get("/domains", domainHandler.GetAll)
	r.Get("/domains/{id}", domainHandler.GetByID)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))
		r.Post("/domains", domainHandler.Create)
		r.Put("/domains/{id}", domainHandler.Update)
		r.Delete("/domains/{id}", domainHandler.Delete)
		r.Post("/chat", chathandler.Chat)
	})

	r.Get("/swagger/*", httpSwagger.Handler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/login.html")
	})
	r.Handle("/*", http.FileServer(http.Dir("./frontend")))
}
