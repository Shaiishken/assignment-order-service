package repository

import (
	"context"
	"database/sql"
	"time"

	"assignment/internal/domain"

	_ "github.com/lib/pq"
)

type OrderPostgresRepository struct {
	db *sql.DB
}

func NewOrderPostgresRepository() (*OrderPostgresRepository, error) {

	connStr := "host=localhost port=5432 user=postgres password=369Zaq57 dbname=AP2 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &OrderPostgresRepository{db: db}, nil
}

func (r *OrderPostgresRepository) Create(ctx context.Context, order *domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO orders (id, customer_id, item_name, amount, status, created_at) VALUES ($1,$2,$3,$4,$5,$6)",
		order.ID, order.CustomerID, order.ItemName, order.Amount, order.Status, order.CreatedAt,
	)
	return err
}

func (r *OrderPostgresRepository) Update(ctx context.Context, order *domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE orders SET status=$1 WHERE id=$2",
		order.Status, order.ID,
	)
	return err
}

func (r *OrderPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {

	row := r.db.QueryRowContext(ctx,
		"SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id=$1",
		id,
	)

	var o domain.Order
	var created time.Time

	err := row.Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &created)
	if err != nil {
		return nil, nil
	}

	o.CreatedAt = created
	return &o, nil
}

func (r *OrderPostgresRepository) GetRevenue(ctx context.Context, customerID string) (map[string]interface{}, error) {

	row := r.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount),0), COUNT(*) FROM orders WHERE customer_id=$1 AND status='paid'",
		customerID,
	)

	var total int64
	var count int

	err := row.Scan(&total, &count)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"customer_id":  customerID,
		"total_amount": total,
		"orders_count": count,
	}

	return result, nil
}
