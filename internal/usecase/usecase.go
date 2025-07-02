package usecase

import (
	"context"
	"ct-backend-course-baonguyen/internal/entity"
	"ct-backend-course-baonguyen/pkg/auth"
	"ct-backend-course-baonguyen/pkg/hashpass"
	"errors"
	"fmt"
	"io"
	_ "net/http"
	"time"
)

type UserStore interface {
	Create(info entity.UserInfo) error
	Query(username string) (entity.UserInfo, error)
	ChangePassword(username string, newPassword string) error
}

type ImageStore interface {
	Store(username string, imagePath string, fileName string) error
	Query(username string) ([]entity.ImageInfo, error)
}

type ImageBucket interface {
	SaveImage(ctx context.Context, name string, r io.Reader) (string, error)
}

func NewUseCase(imageStore ImageStore, userStore UserStore, imageBucket ImageBucket) *ucImplement {
	return &ucImplement{
		userStore:  userStore,
		imageStore: imageStore,
		imgBucket:  imageBucket,
	}
}

type ucImplement struct {
	userStore  UserStore
	imageStore ImageStore
	imgBucket  ImageBucket
}

func (uc *ucImplement) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	if err := uc.userStore.Create(entity.UserInfo{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Address:  req.Address,
	}); err != nil {
		return nil, err
	}

	return &entity.RegisterResponse{UserId: req.Username}, nil
}

func (uc *ucImplement) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := uc.userStore.Query(req.Username)
	if err != nil {
		return nil, err
	}

	// Check if the password matches by comparing with stored hash
	hashedInputPassword := hashpass.HashPasswordLogin(req.Password, user.HashPass)
	if hashedInputPassword != user.HashPass {
		return nil, ErrPasswordMisMatch
	}

	token, err := auth.GenerateToken(user.Username, 24*time.Minute)
	if err != nil {
		return nil, err
	}

	resp := entity.LoginResponse{Token: token}
	return &resp, nil
}

func (uc *ucImplement) Self(ctx context.Context, req *entity.SelfRequest) (*entity.SelfResponse, error) {
	user, err := uc.userStore.Query(req.Username)
	if err != nil {
		fmt.Print("error self uc")
		return nil, err
	}
	selfResp := entity.SelfResponse{
		Username: user.Username,
		Password: user.Password,
		FullName: user.FullName,
		Address:  user.Address,
	}
	return &selfResp, nil
}

func (uc *ucImplement) UploadImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error) {
	imageName, err := uc.imgBucket.SaveImage(ctx, req.FileName, req.Content)
	if err != nil {
		return nil, err
	}

	username := ctx.Value("username").(string)

	// Save image information to storage
	if err := uc.imageStore.Store(username, imageName, req.FileName); err != nil {
		return nil, fmt.Errorf("failed to store image metadata: %w", err)
	}

	return &entity.UploadImageResponse{ImageUrl: imageName}, nil
}

func (uc *ucImplement) ChangePassword(ctx context.Context, req *entity.ChangePasswordRequest) (*entity.ChangePasswordResponse, error) {
	username := ctx.Value("username").(string)

	// Change the password
	err := uc.userStore.ChangePassword(username, req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	return &entity.ChangePasswordResponse{Message: "Password changed"}, nil
}

var ErrPasswordMisMatch = errors.New("password mismatch")
