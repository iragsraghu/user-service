package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/iragsraghu/user-service/internal/models"
	"gorm.io/gorm"
)

var ErrEmailExists = errors.New("email already exists")

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	err := r.DB.WithContext(ctx).Create(u).Error
	if err != nil {
		// MySQL duplicate entry check
		if strings.Contains(err.Error(), "Duplicate entry") {
			return ErrEmailExists
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := r.DB.WithContext(ctx).First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
