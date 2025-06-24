package postgres

import (
	"ct-backend-course-baonguyen/internal/entity"
	"fmt"

	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Store(username string, imagePath string, fileName string) error {
	return fmt.Errorf("not implemented")
}

func (r *ImageRepository) Query(username string) ([]entity.ImageInfo, error) {
	return nil, fmt.Errorf("not implemented")
}
