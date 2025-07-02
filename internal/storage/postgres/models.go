package postgres

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;type:text;not null"`
	HashPass string `gorm:"type:text;not null"`
	FullName string `gorm:"type:text;not null"`
	Address  string `gorm:"type:text"`
}

type Image struct {
	gorm.Model
	UserID uint   `gorm:"index;not null"`
	User   User   `gorm:"foreignKey:UserID"`
	Name   string `gorm:"type:text;not null"`
	Path   string `gorm:"type:text;not null"`
}
