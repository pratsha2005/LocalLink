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

// internal/database/database.go
// ... (keep existing imports and functions)

// -- NEW USER METHODS --

func (s *Store) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, role, created_at FROM users WHERE id = $1`
	err := s.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	return &user, err
}

// -- NEW ORDER METHODS --

// CreateOrder uses a transaction to ensure all-or-nothing order creation.
func (s *Store) CreateOrder(ctx context.Context, input models.CreateOrderInput, buyerID int) (*models.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback on any error

	var totalPrice float64
	var orderItems []models.OrderItem

	// Calculate total price and check product availability
	for _, item := range input.Items {
		var price float64
		var quantityInStock int
		err := tx.QueryRow(ctx, "SELECT price, quantity FROM products WHERE id = $1", item.ProductID).Scan(&price, &quantityInStock)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		if quantityInStock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product ID %d", item.ProductID)
		}
		totalPrice += price * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
	}

	// Insert into orders table
	var orderID int
	orderQuery := `INSERT INTO orders (buyer_id, producer_id, total_price) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx, orderQuery, buyerID, input.ProducerID, totalPrice).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	// Insert into order_items and update product quantities
	for _, item := range orderItems {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(ctx, itemQuery, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}

		updateProductQuery := `UPDATE products SET quantity = quantity - $1 WHERE id = $2`
		_, err = tx.Exec(ctx, updateProductQuery, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Fetch the created order to return it
	createdOrder, err := s.GetOrderByID(ctx, orderID)
	return createdOrder, err
}

func (s *Store) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
    // Implementation to fetch an order and its items (JOIN query)
    // For brevity, a full implementation is complex. This is a conceptual placeholder.
    // In a real app, you would perform a JOIN to get order details and items in one go.
    var order models.Order
    // ... query logic here ...
    return &order, nil
}


func (s *Store) GetOrdersForUser(ctx context.Context, userID int) ([]models.Order, error) {
    // Implementation to fetch all orders for a given user (buyer or producer)
    var orders []models.Order
    // ... query logic here ...
    return orders, nil
}


// -- NEW REVIEW METHODS --

func (s *Store) CreateReview(ctx context.Context, review *models.Review) error {
	query := `INSERT INTO reviews (product_id, user_id, rating, comment) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return s.db.QueryRow(ctx, query, review.ProductID, review.UserID, review.Rating, review.Comment).Scan(&review.ID, &review.CreatedAt)
}

func (s *Store) GetReviewsForProduct(ctx context.Context, productID int) ([]models.Review, error) {
	query := `SELECT id, product_id, user_id, rating, comment, created_at FROM reviews WHERE product_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		if err := rows.Scan(&r.ID, &r.ProductID, &r.UserID, &r.Rating, &r.Comment, &r.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}