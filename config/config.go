package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           string
	MongoURI       string
	MongoDB        string
	MongoCollImage string
	MongoCollUser  string
	// PostgreSQL config
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	PostgresSSLMode  string
	GoogleCredFile   string
	GoogleBucketName string
}

func LoadConfig() Config {
	return Config{
		Port:           GetConfig("PORT"),
		MongoURI:       GetConfig("MONGO_URI"),
		MongoDB:        GetConfig("MONGO_DB"),
		MongoCollImage: GetConfig("MONGO_COLL_IMAGE"),
		MongoCollUser:  GetConfig("MONGO_COLL_USER"),
		// PostgreSQL config with defaults
		PostgresHost:     GetConfig("POSTGRES_HOST"),
		PostgresPort:     GetConfig("POSTGRES_PORT"),
		PostgresUser:     GetConfig("POSTGRES_USER"),
		PostgresPassword: GetConfig("POSTGRES_PASSWORD"),
		PostgresDBName:   GetConfig("POSTGRES_DB"),
		PostgresSSLMode:  GetConfig("POSTGRES_SSLMODE"),
		GoogleCredFile:   GetConfig("GOOGLE_APPLICATION_CREDENTIALS"),
		GoogleBucketName: GetConfig("GOOGLE_APPLICATION_BUCKET"),
	}
}

func GetConfig(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		panic(fmt.Sprintf("Key %s cannot empty", key))
	}
	return val
}

// PostgresConnString returns a formatted PostgreSQL connection string
func (c *Config) PostgresConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword,
		c.PostgresDBName, c.PostgresSSLMode)
}

// GetConfigWithDefault returns environment variable value or default if empty
func GetConfigWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

// ValidateGCSConfig validates Google Cloud Storage configuration
func (c *Config) ValidateGCSConfig() error {
	if c.GoogleCredFile == "" || c.GoogleCredFile == "path/to/your/google/credentials.json" {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is required for GCS integration")
	}

	if c.GoogleBucketName == "" || c.GoogleBucketName == "your-google-bucket-name" {
		return fmt.Errorf("GOOGLE_APPLICATION_BUCKET is required for GCS integration")
	}

	// Check if credentials file exists
	if _, err := os.Stat(c.GoogleCredFile); os.IsNotExist(err) {
		return fmt.Errorf("Google credentials file not found: %s", c.GoogleCredFile)
	}

	return nil
}
