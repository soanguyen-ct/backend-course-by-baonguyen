package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             string
	MongoURI         string
	MongoDB          string
	MongoCollImage   string
	MongoCollUser    string
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
		Port:             GetConfig("PORT"),
		MongoURI:         GetConfig("MONGO_URI"),
		MongoDB:          GetConfig("MONGO_DB"),
		MongoCollImage:   GetConfig("MONGO_COLL_IMAGE"),
		MongoCollUser:    GetConfig("MONGO_COLL_USER"),
		// PostgreSQL config with defaults
		PostgresHost:     GetConfigWithDefault("POSTGRES_HOST", "localhost"),
		PostgresPort:     GetConfigWithDefault("POSTGRES_PORT", "5432"),
		PostgresUser:     GetConfigWithDefault("POSTGRES_USER", "postgres"),
		PostgresPassword: GetConfigWithDefault("POSTGRES_PASSWORD", "postgres"),
		PostgresDBName:   GetConfigWithDefault("POSTGRES_DB", "ct_backend_course"),
		PostgresSSLMode:  GetConfigWithDefault("POSTGRES_SSLMODE", "disable"),
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

func GetConfigWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

// PostgresConnString returns a formatted PostgreSQL connection string
func (c *Config) PostgresConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword,
		c.PostgresDBName, c.PostgresSSLMode)
}
