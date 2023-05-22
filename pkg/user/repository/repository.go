package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"vse.com/4IT428/2023/newsletter/pkg/user/models"
)

// UserRepository represents exported methods for interacting with users
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	// FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUserPassword(ctx context.Context, newPassword, email string) error
}

// SQLUserRepository is a repository for users backed by a SQL database
type SQLUserRepository struct {
	DB *sqlx.DB
}

// CreateUser creates a new user
func (r *SQLUserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `INSERT INTO users (id, email, password)
	          VALUES ($1, $2, $3)`

	_, err := r.DB.ExecContext(ctx, query, user.ID, user.Email, user.Password)
	if err != nil {
		zap.Error(err)
		return models.User{}, err
	}

	return user, nil
}

// GetAllUsers returns all users
func (r *SQLUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT * FROM users`

	var users []models.User
	err := r.DB.SelectContext(ctx, &users, query)
	log.Default().Println("users", err)
	if err != nil {
		zap.Error(err)
		return nil, err
	}
	zap.L().Info("users", zap.Any("users", users))
	log.Default().Println("users", users)
	return users, nil
}

func (r *SQLUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE email=$1", email)
	return user, err
}

func (r *SQLUserRepository) UpdateUserPassword(ctx context.Context, newPassword, email string) error {
	var user models.User
	err := r.DB.GetContext(ctx, &user, "UPDATE users SET password=$1 WHERE email=$2", newPassword, email)
	return err
}
