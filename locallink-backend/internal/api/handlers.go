package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/LocalLink/internal/auth"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	"github.com/LocalLink/internal/models"
	"github.com/LocalLink/internal/websocket"

	"github.com/go-chi/chi/v5"
	gwebsocket "github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	store *database.Store
	cfg   *config.Config
	hub   *websocket.Hub
}

func NewHandler(store *database.Store, cfg *config.Config, hub *websocket.Hub) *Handler {
	return &Handler{store: store, cfg: cfg, hub: hub}
}

// WebSocket Handler
var upgrader = gwebsocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWs(hub *websocket.Hub, w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := websocket.NewClient(hub, conn, userID)
	client.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}

// User Handlers
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user := models.User{Name: input.Name, Email: input.Email, PasswordHash: hashedPassword, Role: input.Role}
	if err := h.store.CreateUser(r.Context(), &user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var input models.LoginUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user, err := h.store.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	token, err := auth.GenerateJWT(user.ID, h.cfg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	user, err := h.store.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var input models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user, err := h.store.UpdateUser(r.Context(), userID, input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// Product Handlers
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	producerID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	product.ProducerID = producerID
	if err := h.store.CreateProduct(r.Context(), &product); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}
	respondWithJSON(w, http.StatusCreated, product)
}

func (h *Handler) GetProductsNearby(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	radius, _ := strconv.Atoi(r.URL.Query().Get("radius"))
	if radius == 0 {
		radius = 5000 // default 5km
	}
	products, err := h.store.GetProductsNearby(r.Context(), lat, lon, radius)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	product, err := h.store.GetProductByID(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	if product.ProducerID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to modify this product")
		return
	}
	var input models.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updatedProduct, err := h.store.UpdateProduct(r.Context(), productID, input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	respondWithJSON(w, http.StatusOK, updatedProduct)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	product, err := h.store.GetProductByID(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	if product.ProducerID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this product")
		return
	}
	if err := h.store.DeleteProduct(r.Context(), productID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Order Handlers
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	buyerID, _ := auth.GetUserIDFromContext(r.Context())
	var input models.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	order, err := h.store.CreateOrder(r.Context(), input, buyerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, order)
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	orders, err := h.store.GetOrdersForUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch orders")
		return
	}
	respondWithJSON(w, http.StatusOK, orders)
}

func (h *Handler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	orderID, _ := strconv.Atoi(chi.URLParam(r, "orderID"))
	order, err := h.store.GetOrderByID(r.Context(), orderID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	if order.BuyerID != userID && order.ProducerID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to view this order")
		return
	}
	respondWithJSON(w, http.StatusOK, order)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	orderID, _ := strconv.Atoi(chi.URLParam(r, "orderID"))
	order, err := h.store.GetOrderByID(r.Context(), orderID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	if order.ProducerID != userID {
		respondWithError(w, http.StatusForbidden, "Only the producer can update the order status")
		return
	}
	var input models.UpdateOrderStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updatedOrder, err := h.store.UpdateOrderStatus(r.Context(), orderID, input.Status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	if client, ok := h.hub.Clients[order.BuyerID]; ok {
		msg := fmt.Sprintf(`{"type": "order_update", "orderId": %d, "status": "%s"}`, orderID, input.Status)
		client.Send <- []byte(msg)
	}
	respondWithJSON(w, http.StatusOK, updatedOrder)
}

// Review Handlers
func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromContext(r.Context())
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	var input models.CreateReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	review := models.Review{ProductID: productID, UserID: userID, Rating: input.Rating, Comment: input.Comment}
	if err := h.store.CreateReview(r.Context(), &review); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create review")
		return
	}
	respondWithJSON(w, http.StatusCreated, review)
}

func (h *Handler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	reviews, err := h.store.GetReviewsForProduct(r.Context(), productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch reviews")
		return
	}
	respondWithJSON(w, http.StatusOK, reviews)
}

// JSON response helpers
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}