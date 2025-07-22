package controller

import (
	"context"
	"ct-backend-course-baonguyen/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockUseCase is a mock implementation of the UseCase interface
type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RegisterResponse), args.Error(1)
}

func (m *MockUseCase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LoginResponse), args.Error(1)
}

func (m *MockUseCase) Self(ctx context.Context, req *entity.SelfRequest) (*entity.SelfResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.SelfResponse), args.Error(1)
}

func (m *MockUseCase) UploadImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UploadImageResponse), args.Error(1)
}

func (m *MockUseCase) ChangePassword(ctx context.Context, req *entity.ChangePasswordRequest) (*entity.ChangePasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ChangePasswordResponse), args.Error(1)
}
