package repositories

import (
	"context"
	"database/sql"

	"boilerplate/server/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type IUserRepository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	List(ctx context.Context) (models.UserSlice, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	if err := u.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context) (models.UserSlice, error) {
	return models.Users().All(ctx, r.db)
}
