package postgres

import (
	"context"
	"ct-backend-course-baonguyen/internal/entity"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ImageRepository struct {
	db      *gorm.DB
	timeout time.Duration
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		db:      db,
		timeout: 3 * time.Second,
	}
}

func (r *ImageRepository) Store(username string, imagePath string, fileName string) error {
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

	image := Image{
		UserID: user.ID,
		Name:   fileName,
		Path:   imagePath,
	}

	result = r.db.WithContext(ctx).Create(&image)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ImageRepository) Query(username string) ([]entity.ImageInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	var images []Image
	result = r.db.WithContext(ctx).Where("user_id = ?", user.ID).Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}

	imageInfos := make([]entity.ImageInfo, len(images))
	for i, img := range images {
		imageInfos[i] = entity.ImageInfo{
			UserName:  username,
			FileName:  img.Name,
			ImagePath: img.Path,
		}
	}

	return imageInfos, nil
}
