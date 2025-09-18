// internal/auth/auth.go
package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/LocalLink/internal/config"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UserIDKey contextKey = "userID"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID int, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"userID":     userID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// AuthMiddleware decodes the share session and packs the session into context
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Could not find bearer token in Authorization header", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIDFloat, ok := claims["userID"].(float64); ok {
					userID := int(userIDFloat)
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, errors.New("no userID found in context")
	}
	return userID, nil
}