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
