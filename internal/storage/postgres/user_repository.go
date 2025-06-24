package postgres

import (
	"ct-backend-course-baonguyen/internal/entity"
	"ct-backend-course-baonguyen/pkg/hashpass"
	"errors"
	"time"

	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db      *gorm.DB
	timeout time.Duration
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:      db,
		timeout: 3 * time.Second,
	}
}

func (r *UserRepository) Create(info entity.UserInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	user := User{
		Username: info.Username,
		HashPass: hashpass.HashPassword(info.Password),
		FullName: info.FullName,
		Address:  info.Address,
	}

	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) Query(username string) (entity.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.UserInfo{}, errors.New("user not found")
		}
		return entity.UserInfo{}, result.Error
	}

	return entity.UserInfo{
		Username: user.Username,
		HashPass: user.HashPass,
		FullName: user.FullName,
		Address:  user.Address,
	}, nil
}

func (r *UserRepository) ChangePassword(username string, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return result.Error
	}

	if hashpass.HashPasswordLogin(newPassword, user.HashPass) == user.HashPass {
		return errors.New("same password")
	}

	newHashedPassword := hashpass.HashPasswordLogin(newPassword, user.HashPass)
	result = r.db.Model(&user).Update("hash_pass", newHashedPassword)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
