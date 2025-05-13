package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Created  time.Time `json:"created"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created`

	args := []any{user.Username, user.Email, user.Password}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Created)
}

func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, sql.ErrNoRows
	}

	query := `
		SELECT id, username, email, password, created
		FROM users
		WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Created,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3
		WHERE id = $4`

	args := []any{
		user.Username,
		user.Email,
		user.Password,
		user.ID,
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

func (m UserModel) Authenticate(user *User) (*User, error) {
	query := `
		SELECT id, username, password, created
		FROM users
		WHERE email = $1`

	args := []any{
		user.Email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userID int
	var username string
	var password string
	var created time.Time

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&userID,
		&username,
		&password,
		&created,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email")
		} else {
			return nil, err
		}
	}

	if user.Password != password {
		return nil, errors.New("invalid password")
	}

	user.ID = userID
	user.Username = username
	user.Password = password
	user.Created = created

	return user, nil
}
