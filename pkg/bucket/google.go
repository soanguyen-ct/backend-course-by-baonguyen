package bucket

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
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

func (c *GoogleStorageClient) SaveImage(ctx context.Context, fileName string, fileReader io.Reader) (string, error) {
	// TODO
	return "https://storage.googleapis.com" + "<input yourself>", nil
}
