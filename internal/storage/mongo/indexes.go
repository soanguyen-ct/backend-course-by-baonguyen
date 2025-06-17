package mongostore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUserCollectionIndex(userColl *userCollection) error {
	c := userColl.client

	// TODO 2: Create unique index for 'username'

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := c.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	return nil
}
