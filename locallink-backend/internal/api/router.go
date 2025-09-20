package api

import (
	"net/http"

	"github.com/LocalLink/internal/auth"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	"github.com/LocalLink/internal/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(store *database.Store, cfg *config.Config, hub *websocket.Hub) *chi.Mux {
	r := chi.NewRouter()
	h := NewHandler(store, cfg, hub)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public Routes
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)
	r.Get("/products/nearby", h.GetProductsNearby)
	r.Get("/products/{productID}/reviews", h.GetProductReviews)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(cfg))

		// WebSocket connection
		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			ServeWs(hub, w, r)
		})

		// User & Profile Management
		r.Get("/users/me", h.GetUserProfile)
		r.Put("/users/me", h.UpdateUserProfile)

		// Product Management
		r.Post("/products", h.CreateProduct)
		r.Put("/products/{productID}", h.UpdateProduct)
		r.Delete("/products/{productID}", h.DeleteProduct)

		// Order Management
		r.Post("/orders", h.CreateOrder)
		r.Get("/orders", h.GetUserOrders)
		r.Get("/orders/{orderID}", h.GetOrderDetails)
		r.Put("/orders/{orderID}/status", h.UpdateOrderStatus)

		// Review Management
		r.Post("/products/{productID}/reviews", h.CreateReview)
	})

	return r
}