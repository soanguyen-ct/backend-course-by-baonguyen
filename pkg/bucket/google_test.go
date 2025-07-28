package bucket

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGoogleStorageClient_generateUniqueFileName(t *testing.T) {
	client := &GoogleStorageClient{
		Bucket: "test-bucket",
	}

	tests := []struct {
		name        string
		fileName    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid jpg file",
			fileName: "test.jpg",
			wantErr:  false,
		},
		{
			name:     "valid png file",
			fileName: "test.png",
			wantErr:  false,
		},
		{
			name:     "valid jpeg file",
			fileName: "test.jpeg",
			wantErr:  false,
		},
		{
			name:     "valid gif file",
			fileName: "test.gif",
			wantErr:  false,
		},
		{
			name:     "valid webp file",
			fileName: "test.webp",
			wantErr:  false,
		},
		{
			name:        "invalid txt file",
			fileName:    "test.txt",
			wantErr:     true,
			errContains: "unsupported file type",
		},
		{
			name:        "invalid exe file",
			fileName:    "test.exe",
			wantErr:     true,
			errContains: "unsupported file type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.generateUniqueFileName(tt.fileName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("generateUniqueFileName() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("generateUniqueFileName() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("generateUniqueFileName() unexpected error = %v", err)
				return
			}

			// Check that result contains date path and UUID
			if !strings.Contains(result, "/") {
				t.Errorf("generateUniqueFileName() result should contain date path separators")
			}

			// Check that result ends with original extension
			ext := strings.ToLower(filepath.Ext(tt.fileName))
			if !strings.HasSuffix(strings.ToLower(result), ext) {
				t.Errorf("generateUniqueFileName() result should end with extension %s, got %s", ext, result)
			}

			// Check that result contains current year
			currentYear := time.Now().Year()
			yearStr := fmt.Sprintf("%04d", currentYear)
			if !strings.Contains(result, yearStr) {
				t.Errorf("generateUniqueFileName() result should contain current year %s", yearStr)
			}
		})
	}
}

func TestGoogleStorageClient_detectContentType(t *testing.T) {
	client := &GoogleStorageClient{}

	tests := []struct {
		fileName string
		expected string
	}{
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.png", "image/png"},
		{"test.gif", "image/gif"},
		{"test.webp", "image/webp"},
		{"test.unknown", "application/octet-stream"},
		{"test.JPG", "image/jpeg"}, // Test case insensitivity
		{"test.PNG", "image/png"},  // Test case insensitivity
	}

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			result := client.detectContentType(tt.fileName)
			if result != tt.expected {
				t.Errorf("detectContentType(%s) = %s, want %s", tt.fileName, result, tt.expected)
			}
		})
	}
}

func TestGoogleStorageClient_isValidImageExtension(t *testing.T) {
	client := &GoogleStorageClient{}

	tests := []struct {
		extension string
		expected  bool
	}{
		{".jpg", true},
		{".jpeg", true},
		{".png", true},
		{".gif", true},
		{".webp", true},
		{".txt", false},
		{".exe", false},
		{".pdf", false},
		{".doc", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.extension, func(t *testing.T) {
			result := client.isValidImageExtension(tt.extension)
			if result != tt.expected {
				t.Errorf("isValidImageExtension(%s) = %v, want %v", tt.extension, result, tt.expected)
			}
		})
	}
}
