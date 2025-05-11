package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Price     int64     `json:"price"`
	Stock     int64     `json:"stock"`
	ISBN      string    `json:"isbn"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"-"`
}

type BookModel struct {
	DB *sql.DB
}

func (m BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books (title, author, price, stock, isbn, image)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	args := []any{book.Title, book.Author, book.Price, book.Stock, book.ISBN, book.Image}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt)
}

func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT id, title, author, price, stock, isbn, image, created_at
		FROM books
		WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Price,
		&book.Stock,
		&book.ISBN,
		&book.Image,
		&book.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, price = $3, stock = $4, isbn = $5, image = $6
		WHERE id = $7`

	args := []any{
		book.Title,
		book.Author,
		book.Price,
		book.Stock,
		book.ISBN,
		book.Image,
		book.ID,
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

func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return sql.ErrNoRows
	}

	query := `
		DELETE FROM books
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

func (m BookModel) GetAll() ([]*Book, error) {
	query := `
		SELECT id, title, author, price, stock, isbn, image, created_at
		FROM books
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Price,
			&book.Stock,
			&book.ISBN,
			&book.Image,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
