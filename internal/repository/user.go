package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// users table
	User struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		Password  []byte    `db:"password"`
		UpdatedAt time.Time `db:"updated_at"`
		CreatedAt time.Time `db:"created_at"`
	}

	CreateUserParams struct {
		Name     string
		Password string
	}

	UpdateNameParams struct {
		ID   uuid.UUID
		Name string
	}

	UpdatePassParams struct {
		ID       uuid.UUID
		Password string
	}
)

func (r *Repository) CreateUser(ctx context.Context, params CreateUserParams) (uuid.UUID, error) {
	userID := uuid.New()
	hased, err := hashPassword(params.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, "INSERT INTO users (id, name, password) VALUES (?, ?, ?)", userID, params.Name, hased); err != nil {
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetUserID(ctx context.Context, name string) (uuid.UUID, error) {
	user := &User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE name = ?", name); err != nil {
		return uuid.Nil, fmt.Errorf("select user: %w", err)
	}

	return user.ID, nil
}

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id = ?", userID); err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}

	return user, nil
}

func (r *Repository) UpdateName(ctx context.Context, params UpdateNameParams) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE users SET name = ? WHERE id = ?", params.Name, params.ID); err != nil {
		return fmt.Errorf("update user name: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePass(ctx context.Context, params UpdatePassParams) error {

	hashed, err := hashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, "UPDATE users SET password = ? WHERE id = ?", hashed, params.ID); err != nil {
		return fmt.Errorf("update user password: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (r *Repository) CheckPass(ctx context.Context, userID uuid.UUID, password string) (bool, error) {
	user := &User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id = ?", userID); err != nil {
		return false, fmt.Errorf("select user: %w", err)
	}

	return comparePassword(password, user.Password), nil
}

func hashPassword(pass string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return hash, err
}

func comparePassword(pass string, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(pass)) == nil
}
