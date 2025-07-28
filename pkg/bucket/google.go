package bucket

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
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
	log.Printf("[GCS] Starting upload for file: %s", fileName)

	// Generate unique filename to prevent conflicts
	uniqueFileName, err := c.generateUniqueFileName(fileName)
	if err != nil {
		log.Printf("[GCS] Failed to generate unique filename for %s: %v", fileName, err)
		return "", fmt.Errorf("failed to generate unique filename: %w", err)
	}

	log.Printf("[GCS] Generated unique filename: %s", uniqueFileName)

	// Wrap reader with size limit to prevent large file uploads
	limitedReader := io.LimitReader(fileReader, MaxFileSize+1)

	// Create object handle
	obj := c.StorageClient.Bucket(c.Bucket).Object(uniqueFileName)

	// Create writer with context and timeout
	writer := obj.NewWriter(ctx)
	defer func() {
		if closeErr := writer.Close(); closeErr != nil && err == nil {
			log.Printf("[GCS] Error closing writer for %s: %v", uniqueFileName, closeErr)
			err = fmt.Errorf("failed to finalize upload: %w", closeErr)
		}
	}()

	// Set content type based on file extension
	contentType := c.detectContentType(fileName)
	if contentType != "" {
		writer.ContentType = contentType
		log.Printf("[GCS] Set content type: %s for %s", contentType, uniqueFileName)
	}

	// Set cache control for better performance
	writer.CacheControl = "public, max-age=86400" // Cache for 24 hours

	// Copy file data to GCS with size validation
	log.Printf("[GCS] Starting file transfer to bucket: %s", c.Bucket)
	bytesWritten, err := io.Copy(writer, limitedReader)
	if err != nil {
		log.Printf("[GCS] Failed to upload %s: %v", uniqueFileName, err)
		return "", fmt.Errorf("failed to upload file to GCS: %w", err)
	}

	// Check if file size exceeds limit
	if bytesWritten > MaxFileSize {
		log.Printf("[GCS] File %s exceeds size limit: %d bytes", uniqueFileName, bytesWritten)
		return "", fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxFileSize)
	}

	log.Printf("[GCS] Successfully uploaded %d bytes for %s", bytesWritten, uniqueFileName)

	// Close writer to finalize upload
	if err := writer.Close(); err != nil {
		log.Printf("[GCS] Failed to finalize upload for %s: %v", uniqueFileName, err)
		return "", fmt.Errorf("failed to finalize upload: %w", err)
	}

	// Return public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.Bucket, uniqueFileName)
	log.Printf("[GCS] Upload completed successfully. Public URL: %s", publicURL)
	return publicURL, nil
}

// generateUniqueFileName creates a unique filename using UUID while preserving the original extension
func (c *GoogleStorageClient) generateUniqueFileName(originalFileName string) (string, error) {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(originalFileName))

	// Validate file extension (basic image types only)
	if !c.isValidImageExtension(ext) {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	// Generate UUID
	uniqueID := uuid.New().String()

	// Create date-based folder structure: YYYY/MM/DD/
	now := time.Now()
	datePath := fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day())

	// Combine date path with unique filename
	uniqueFileName := fmt.Sprintf("%s/%s%s", datePath, uniqueID, ext)

	return uniqueFileName, nil
}

// detectContentType determines the MIME type based on file extension
func (c *GoogleStorageClient) detectContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	contentType := mime.TypeByExtension(ext)

	// If mime package doesn't recognize it, set common image types manually
	if contentType == "" {
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

	return contentType
}

// isValidImageExtension validates that the file extension is an allowed image type
func (c *GoogleStorageClient) isValidImageExtension(ext string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}

	return false
}
