package controller

import (
	"bytes"
	"ct-backend-course-baonguyen/internal/entity"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test data constants
const (
	TestUsername  = "testuser"
	TestPassword  = "password123"
	ShortPassword = "short"
	TestFullName  = "Test User"
	TestAddress   = "123 Test Street"
	TestJWTToken  = "mock.jwt.token"
	TestImageURL  = "/static/images/test.jpg"
)

// setupTestEcho creates a test Echo instance with routes configured
func setupTestEcho(uc UseCase) *echo.Echo {
	e := echo.New()
	// Set up a mock validator that always passes
	e.Validator = &mockValidator{}
	handler := NewHandler(uc)

	public := e.Group("/api/public")
	private := e.Group("/api/private")
	private.Use(mockAuthMiddleware()) // Mock JWT authentication

	// Register routes
	public.POST("/register", handler.Register)
	public.POST("/login", handler.Login)
	private.GET("/self", handler.Self)
	private.POST("/upload", handler.UploadImage)
	private.PUT("/change-password", handler.ChangePassword)

	return e
}

// mockValidator is a simple validator that always passes validation
type mockValidator struct{}

func (mv *mockValidator) Validate(i interface{}) error {
	return nil // Always pass validation in tests
}

// mockAuthMiddleware creates a mock authentication middleware for testing
func mockAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set a test username in context to simulate authenticated user
			c.Set("username", TestUsername)
			return next(c)
		}
	}
}

// Helper functions
func createJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func executeRequest(e *echo.Echo, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func parseJSONResponse(rec *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(rec.Body.Bytes(), target)
}

func createMultipartRequest(url string, filename string, content []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form file
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, err
	}

	// Write file content
	if _, err := part.Write(content); err != nil {
		return nil, err
	}

	// Close the writer to finalize the multipart content
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func createTestImageFile() []byte {
	// Create a minimal test image (simple PNG header + data)
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, // IHDR chunk length
		0x49, 0x48, 0x44, 0x52, // IHDR
		0x00, 0x00, 0x00, 0x01, // Width: 1
		0x00, 0x00, 0x00, 0x01, // Height: 1
		0x08, 0x02, 0x00, 0x00, 0x00, // Bit depth, Color type, etc.
		0x90, 0x77, 0x53, 0xDE, // CRC
		0x00, 0x00, 0x00, 0x0C, // IDAT chunk length
		0x49, 0x44, 0x41, 0x54, // IDAT
		0x08, 0x99, 0x01, 0x01, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01,
		0xE2, 0x21, 0xBC, 0x33, // CRC
		0x00, 0x00, 0x00, 0x00, // IEND chunk length
		0x49, 0x45, 0x4E, 0x44, // IEND
		0xAE, 0x42, 0x60, 0x82, // CRC
	}
}

// Test fixtures
func validRegisterRequest() entity.RegisterRequest {
	return entity.RegisterRequest{
		Username: TestUsername,
		Password: TestPassword,
		FullName: TestFullName,
		Address:  TestAddress,
	}
}

func validLoginRequest() entity.LoginRequest {
	return entity.LoginRequest{
		Username: TestUsername,
		Password: TestPassword,
	}
}

func validChangePasswordRequest() entity.ChangePasswordRequest {
	return entity.ChangePasswordRequest{
		OldPassword: TestPassword,
		NewPassword: "newpassword123",
	}
}

// TestHandler_Register tests the Register endpoint
func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "Valid registration success",
			requestBody: validRegisterRequest(),
			setupMock: func(uc *MockUseCase) {
				uc.On("Register", mock.Anything, mock.AnythingOfType("*entity.RegisterRequest")).
					Return(&entity.RegisterResponse{UserId: TestUsername}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "UseCase error - duplicate user",
			requestBody: validRegisterRequest(),
			setupMock: func(uc *MockUseCase) {
				uc.On("Register", mock.Anything, mock.AnythingOfType("*entity.RegisterRequest")).
					Return(nil, errors.New("user already exists"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUC := new(MockUseCase)
			tt.setupMock(mockUC)

			e := setupTestEcho(mockUC)

			req, err := createJSONRequest("POST", "/api/public/register", tt.requestBody)
			assert.NoError(t, err)

			// Execute
			rec := executeRequest(e, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.RegisterResponse
				err := parseJSONResponse(rec, &response)
				assert.NoError(t, err)
				assert.Equal(t, TestUsername, response.UserId)
			}

			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

// TestHandler_Login tests the Login endpoint
func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "Valid login success",
			requestBody: validLoginRequest(),
			setupMock: func(uc *MockUseCase) {
				uc.On("Login", mock.Anything, mock.AnythingOfType("*entity.LoginRequest")).
					Return(&entity.LoginResponse{Token: TestJWTToken}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Invalid credentials",
			requestBody: validLoginRequest(),
			setupMock: func(uc *MockUseCase) {
				uc.On("Login", mock.Anything, mock.AnythingOfType("*entity.LoginRequest")).
					Return(nil, errors.New("password mismatch"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUC := new(MockUseCase)
			tt.setupMock(mockUC)

			e := setupTestEcho(mockUC)

			req, err := createJSONRequest("POST", "/api/public/login", tt.requestBody)
			assert.NoError(t, err)

			// Execute
			rec := executeRequest(e, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.LoginResponse
				err := parseJSONResponse(rec, &response)
				assert.NoError(t, err)
				assert.Equal(t, TestJWTToken, response.Token)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

// TestHandler_Self tests the Self endpoint
func TestHandler_Self(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid authenticated request",
			setupMock: func(uc *MockUseCase) {
				uc.On("Self", mock.Anything, mock.AnythingOfType("*entity.SelfRequest")).
					Return(&entity.SelfResponse{
						Username: TestUsername,
						Password: TestPassword,
						FullName: TestFullName,
						Address:  TestAddress,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User not found",
			setupMock: func(uc *MockUseCase) {
				uc.On("Self", mock.Anything, mock.AnythingOfType("*entity.SelfRequest")).
					Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUC := new(MockUseCase)
			tt.setupMock(mockUC)

			e := setupTestEcho(mockUC)

			req := httptest.NewRequest("GET", "/api/private/self", nil)

			// Execute
			rec := executeRequest(e, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.SelfResponse
				err := parseJSONResponse(rec, &response)
				assert.NoError(t, err)
				assert.Equal(t, TestUsername, response.Username)
			}

			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

// TestHandler_UploadImage tests the UploadImage endpoint
func TestHandler_UploadImage(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() (*http.Request, error)
		setupMock      func(*MockUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid image upload success",
			setupRequest: func() (*http.Request, error) {
				imageData := createTestImageFile()
				return createMultipartRequest("/api/private/upload", "test.jpg", imageData)
			},
			setupMock: func(uc *MockUseCase) {
				uc.On("UploadImage", mock.Anything, mock.AnythingOfType("*entity.UploadImageRequest")).
					Return(&entity.UploadImageResponse{ImageUrl: TestImageURL}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "UseCase upload failure",
			setupRequest: func() (*http.Request, error) {
				imageData := createTestImageFile()
				return createMultipartRequest("/api/private/upload", "test.jpg", imageData)
			},
			setupMock: func(uc *MockUseCase) {
				uc.On("UploadImage", mock.Anything, mock.AnythingOfType("*entity.UploadImageRequest")).
					Return(nil, errors.New("storage bucket error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUC := new(MockUseCase)
			tt.setupMock(mockUC)

			e := setupTestEcho(mockUC)

			req, err := tt.setupRequest()
			assert.NoError(t, err)

			// Execute
			rec := executeRequest(e, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.UploadImageResponse
				err := parseJSONResponse(rec, &response)
				assert.NoError(t, err)
				assert.Equal(t, TestImageURL, response.ImageUrl)
			}

			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

// TestHandler_ChangePassword tests the ChangePassword endpoint
func TestHandler_ChangePassword(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "Valid password change success",
			requestBody: validChangePasswordRequest(),
			setupMock: func(uc *MockUseCase) {
				uc.On("ChangePassword", mock.Anything, mock.AnythingOfType("*entity.ChangePasswordRequest")).
					Return(&entity.ChangePasswordResponse{Message: "Password changed"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "UseCase password change failure",
			requestBody: validChangePasswordRequest(),
			setupMock: func(uc *MockUseCase) {
				uc.On("ChangePassword", mock.Anything, mock.AnythingOfType("*entity.ChangePasswordRequest")).
					Return(nil, errors.New("failed to change password"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "failed to change password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUC := new(MockUseCase)
			tt.setupMock(mockUC)

			e := setupTestEcho(mockUC)

			req, err := createJSONRequest("PUT", "/api/private/change-password", tt.requestBody)
			assert.NoError(t, err)

			// Execute
			rec := executeRequest(e, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.ChangePasswordResponse
				err := parseJSONResponse(rec, &response)
				assert.NoError(t, err)
				assert.Equal(t, "Password changed", response.Message)
			}

			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

// TestHandler_MalformedJSON tests handling of malformed JSON requests
func TestHandler_MalformedJSON(t *testing.T) {
	mockUC := new(MockUseCase)
	e := setupTestEcho(mockUC)

	req := httptest.NewRequest("POST", "/api/public/register", strings.NewReader(`{"invalid": json}`))
	req.Header.Set("Content-Type", "application/json")

	rec := executeRequest(e, req)

	// Should return an error status due to malformed JSON
	assert.NotEqual(t, http.StatusOK, rec.Code)
}

// TestHandler_EmptyBody tests handling of empty request bodies
func TestHandler_EmptyBody(t *testing.T) {
	mockUC := new(MockUseCase)
	// Set up mock to handle the unexpected call with empty data
	mockUC.On("Register", mock.Anything, mock.AnythingOfType("*entity.RegisterRequest")).
		Return(nil, errors.New("validation error"))

	e := setupTestEcho(mockUC)

	req := httptest.NewRequest("POST", "/api/public/register", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := executeRequest(e, req)

	// Should return an error status due to empty body
	assert.NotEqual(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}
