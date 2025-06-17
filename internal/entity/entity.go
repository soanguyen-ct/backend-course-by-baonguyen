package entity

import (
	"io"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,gt=0"`
	Password string `json:"password" validate:"required,gt=7"`
}

type LoginResponse struct {
	Token string
}

type SelfRequest struct {
	Username string
}

type SelfResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

type UploadImageRequest struct {
	FileName string
	Content  io.Reader
}

type UploadImageResponse struct {
	ImageUrl string `json:"imageURL"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
	HashPass string
}

type ImageInfo struct {
	UserName  string `json:"username"`
	ImagePath string `json:"image_path"`
	FileName  string `json:"file_name"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}
