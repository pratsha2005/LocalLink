// internal/database/database.go
package database

import (
	"context"
	"fmt"
	"github.com/LocalLink/internal/models"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func Connect(databaseURL string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return dbpool
}

// User methods
func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return s.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	return &user, err
}

// Product methods
func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (producer_id, name, description, price, quantity, location) 
              VALUES ($1, $2, $3, $4, $5, ST_MakePoint($6, $7)::geography) 
              RETURNING id, created_at`
	// Note: PostGIS stores as (longitude, latitude)
	return s.db.QueryRow(ctx, query, product.ProducerID, product.Name, product.Description, product.Price, product.Quantity, product.Longitude, product.Latitude).Scan(&product.ID, &product.CreatedAt)
}

func (s *Store) GetProductsNearby(ctx context.Context, lat, lon float64, radius int) ([]models.Product, error) {
	query := `
        SELECT id, producer_id, name, description, price, quantity, ST_Y(location::geometry), ST_X(location::geometry), created_at
        FROM products
        WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3)`

	rows, err := s.db.Query(ctx, query, lon, lat, radius) // lon, lat for PostGIS
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.ProducerID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.Latitude, &p.Longitude, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}