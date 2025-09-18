package api

import (
	"github.com/LocalLink/internal/auth"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	// "net/http" 
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter sets up and returns the main application router with all API endpoints.
func NewRouter(store *database.Store, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	h := NewHandler(store, cfg)

	// --- Middleware ---
	// Logger provides readable logs for each request.
	// Recoverer handles panics and returns a 500 error.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- Public Routes ---
	// These endpoints do not require authentication.
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)
	r.Get("/products/nearby", h.GetProductsNearby)
	r.Get("/products/{productID}/reviews", h.GetProductReviews)

	// --- Protected Routes ---
	// These endpoints require a valid JWT. The AuthMiddleware handles token validation.
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(cfg))

		// User & Profile Management
		r.Get("/users/me", h.GetUserProfile)
		// TODO: Add PUT /users/me for updating user profiles.

		// Product Management
		r.Post("/products", h.CreateProduct)
		// TODO: Add PUT /products/{productID} and DELETE /products/{productID}.

		// Order Management
		r.Post("/orders", h.CreateOrder)
		r.Get("/orders", h.GetUserOrders)
		// TODO: Add GET /orders/{orderID} and PUT /orders/{orderID}/status.

		// Review Management
		r.Post("/products/{productID}/reviews", h.CreateReview)
	})

	return r
}