// internal/api/router.go
package api

import (
	"github.com/LocalLink/internal/auth"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	// "net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(store *database.Store, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	h := NewHandler(store, cfg)

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)
	r.Get("/products/nearby", h.GetProductsNearby)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(cfg))
		r.Post("/products", h.CreateProduct)
		// Add other protected routes here (e.g., POST /orders)
	})

	return r
}