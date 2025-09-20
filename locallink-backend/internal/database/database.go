package database

import (
	"context"
	"fmt"
	"log"

	"github.com/LocalLink/internal/models"

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
	if err = dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")
	return dbpool
}

// User Methods
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

func (s *Store) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, role, created_at FROM users WHERE id = $1`
	err := s.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	return &user, err
}

func (s *Store) UpdateUser(ctx context.Context, userID int, input models.UpdateUserInput) (*models.User, error) {
	if input.Name != nil {
		query := `UPDATE users SET name = $1 WHERE id = $2`
		_, err := s.db.Exec(ctx, query, *input.Name, userID)
		if err != nil {
			return nil, err
		}
	}
	return s.GetUserByID(ctx, userID)
}

// Product Methods
func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (producer_id, name, description, price, quantity, location) 
              VALUES ($1, $2, $3, $4, $5, ST_MakePoint($6, $7)::geography) RETURNING id, created_at`
	return s.db.QueryRow(ctx, query, product.ProducerID, product.Name, product.Description, product.Price, product.Quantity, product.Longitude, product.Latitude).Scan(&product.ID, &product.CreatedAt)
}

func (s *Store) GetProductsNearby(ctx context.Context, lat, lon float64, radius int) ([]models.Product, error) {
	query := `SELECT id, producer_id, name, description, price, quantity, ST_Y(location::geometry), ST_X(location::geometry), created_at
              FROM products WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3)`
	rows, err := s.db.Query(ctx, query, lon, lat, radius)
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

func (s *Store) GetProductByID(ctx context.Context, productID int) (*models.Product, error) {
	var p models.Product
	query := `SELECT id, producer_id, name, description, price, quantity, ST_Y(location::geometry), ST_X(location::geometry), created_at FROM products WHERE id = $1`
	err := s.db.QueryRow(ctx, query, productID).Scan(&p.ID, &p.ProducerID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.Latitude, &p.Longitude, &p.CreatedAt)
	return &p, err
}

func (s *Store) UpdateProduct(ctx context.Context, productID int, input models.UpdateProductInput) (*models.Product, error) {
	query := `UPDATE products SET name = COALESCE($1, name), description = COALESCE($2, description), price = COALESCE($3, price), quantity = COALESCE($4, quantity) WHERE id = $5`
	_, err := s.db.Exec(ctx, query, input.Name, input.Description, input.Price, input.Quantity, productID)
	if err != nil {
		return nil, err
	}
	return s.GetProductByID(ctx, productID)
}

func (s *Store) DeleteProduct(ctx context.Context, productID int) error {
	_, err := s.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	return err
}

// Order Methods
func (s *Store) CreateOrder(ctx context.Context, input models.CreateOrderInput, buyerID int) (*models.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var totalPrice float64
	var orderItems []models.OrderItem

	for _, item := range input.Items {
		var price float64
		var quantityInStock int
		err := tx.QueryRow(ctx, "SELECT price, quantity FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&price, &quantityInStock)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		if quantityInStock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product ID %d", item.ProductID)
		}
		totalPrice += price * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{ProductID: item.ProductID, Quantity: item.Quantity, Price: price})
	}

	var orderID int
	orderQuery := `INSERT INTO orders (buyer_id, producer_id, total_price) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx, orderQuery, buyerID, input.ProducerID, totalPrice).Scan(&orderID)
	if err != nil {
		return nil, err
	}

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
	return s.GetOrderByID(ctx, orderID)
}

func (s *Store) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	orderQuery := `SELECT id, buyer_id, producer_id, total_price, status, created_at FROM orders WHERE id = $1`
	err := s.db.QueryRow(ctx, orderQuery, orderID).Scan(&order.ID, &order.BuyerID, &order.ProducerID, &order.TotalPrice, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	itemsQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := s.db.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}

func (s *Store) GetOrdersForUser(ctx context.Context, userID int) ([]models.Order, error) {
	query := `SELECT id FROM orders WHERE buyer_id = $1 OR producer_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var orderID int
		if err := rows.Scan(&orderID); err != nil {
			return nil, err
		}
		fullOrder, err := s.GetOrderByID(ctx, orderID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *fullOrder)
	}
	return orders, nil
}

func (s *Store) UpdateOrderStatus(ctx context.Context, orderID int, status string) (*models.Order, error) {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := s.db.Exec(ctx, query, status, orderID)
	if err != nil {
		return nil, err
	}
	return s.GetOrderByID(ctx, orderID)
}

// Review Methods
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