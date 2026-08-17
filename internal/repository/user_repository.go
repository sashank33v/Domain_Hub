package repository

import (
	"database/sql"
	"domainhub/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (r *UserRepository) Create(user *models.User) error {

	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`
	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
