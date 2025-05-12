package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	BookID    int64     `json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderModel struct {
	DB *sql.DB
}

func (m OrderModel) Insert(order *Order) error {
	query := `
		INSERT INTO orders (user_id, book_id)
		VALUES ($1, $2)
		RETURNING id, created_at`

	args := []any{order.UserID, order.BookID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID, &order.CreatedAt)
}

func (m OrderModel) Get(id int64) (*Order, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT id, user_id, book_id, created_at
		FROM orders
		WHERE id = $1`

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.BookID,
		&order.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	return &order, nil
}

func (m OrderModel) Update(order *Order) error {
	query := `
		UPDATE orders
		SET user_id = $1, book_id = $2
		WHERE id = $3`

	args := []any{
		order.UserID,
		order.BookID,
		order.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return sql.ErrNoRows
		default:
			return err
		}
	}

	return nil
}

func (m OrderModel) Delete(id int64) error {
	if id < 1 {
		return sql.ErrNoRows
	}

	query := `
		DELETE FROM orders
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m OrderModel) GetByUserID(userID int64) ([]*Order, error) {
	query := `
		SELECT id, user_id, book_id, created_at
		FROM orders
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, sql.ErrNoRows
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.BookID,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
