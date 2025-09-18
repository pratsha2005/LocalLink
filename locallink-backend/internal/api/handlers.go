// internal/api/handlers.go
package api

import (
	"encoding/json"
	"errors"
	"github.com/LocalLink/internal/auth"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	"github.com/LocalLink/internal/models"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	store *database.Store
	cfg   *config.Config
}

func NewHandler(store *database.Store, cfg *config.Config) *Handler {
	return &Handler{store: store, cfg: cfg}
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

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         input.Role,
	}

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
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	radStr := r.URL.Query().Get("radius")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid 'lat' parameter")
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid 'lon' parameter")
		return
	}
	radius, err := strconv.Atoi(radStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid 'radius' parameter")
		return
	}

	products, err := h.store.GetProductsNearby(r.Context(), lat, lon, radius)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	respondWithJSON(w, http.StatusOK, products)
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