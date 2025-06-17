package mongostore

import (
	"context"
	"ct-backend-course-baonguyen/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewImageCollection(db *mongo.Database, collName string) *imageCollection {
	return &imageCollection{
		client:  db.Collection(collName),
		timeout: 3 * time.Second,
	}
}

type imageCollection struct {
	client  *mongo.Collection
	timeout time.Duration
}

func (c *imageCollection) Save(info entity.ImageInfo) error {
	doc := NewImageDocument(info)

	ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFn()

	_, err := c.client.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}
