package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/fylgushev/go-diplom-final/internal/domain/entities"
	"github.com/fylgushev/go-diplom-final/internal/domain/interfaces"
)

// UserRepository реализует интерфейс UserRepository с использованием GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных
func (r *UserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByLogin получает пользователя по логину
func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("username = ?", login).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (r *UserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser удаляет пользователя по ID
func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	err := r.db.WithContext(ctx).Where("id = ?", userID).Delete(&entities.User{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}