package bucket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/google/uuid"
)

const (
	// MaxFileSize represents the maximum allowed file size (10MB)
	MaxFileSize = 10 << 20 // 10MB in bytes
)

var (
	// Error definitions for better error handling
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrFileSizeExceeded    = errors.New("file size exceeds maximum limit")
	ErrUploadFailed        = errors.New("failed to upload file to GCS")
	ErrUploadFinalize      = errors.New("failed to finalize upload")
	ErrUniqueFileName      = errors.New("failed to generate unique filename")
)

// Simple list of allowed image file extensions
var allowedExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

func MustNewGoogleStorageClient(ctx context.Context, bucketName, credentialsFile string) *GoogleStorageClient {
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		panic(err.Error())
	}

	return &GoogleStorageClient{
		Bucket:        bucketName,
		StorageClient: storageClient,
	}
}

type GoogleStorageClient struct {
	Bucket        string
	StorageClient *storage.Client
}

// SaveImage uploads an image file to Google Cloud Storage and returns the public URL
func (c *GoogleStorageClient) SaveImage(ctx context.Context, fileName string, fileReader io.Reader) (string, error) {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileName))
	if !c.isValidImageExtension(ext) {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedFileType, ext)
	}

	// Create unique filename with date folder
	uniqueID := uuid.New().String()
	now := time.Now()
	year, month, day := now.Date()
	uniqueFileName := fmt.Sprintf("%d/%02d/%02d/%s%s", year, int(month), day, uniqueID, ext)

	// Upload with size limit
	obj := c.StorageClient.Bucket(c.Bucket).Object(uniqueFileName)
	writer := obj.NewWriter(ctx)

	// Set content type
	writer.ContentType = c.detectContentType(fileName)

	// Limit to MaxFileSize
	limitedReader := io.LimitReader(fileReader, MaxFileSize+1)
	bytesWritten, err := io.Copy(writer, limitedReader)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// check if file is too large
	if bytesWritten > MaxFileSize {
		return "", fmt.Errorf("%w: %d bytes", ErrFileSizeExceeded, bytesWritten)
	}

	// complete upload
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadFinalize, err)
	}

	// return URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.Bucket, uniqueFileName)
	return publicURL, nil
}

// detectContentType determines the MIME type based on file extension
func (c *GoogleStorageClient) detectContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	// Simple switch for common image types
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// isValidImageExtension validates that the file extension is an allowed image type
func (c *GoogleStorageClient) isValidImageExtension(ext string) bool {
	// Simple loop through allowed extensions
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
