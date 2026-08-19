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
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
	Select id, email, password_hash, role, created_at
	From users
	Where email = $1;
	`
	var user models.User

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
