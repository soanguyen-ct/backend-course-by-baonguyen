package mongostore

import (
	"context"
	"ct-backend-course-baonguyen/internal/entity"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestImageCollection_Save(t *testing.T) {
	db := MustDatabase(
		"mongodb+srv://bao:baopass@cluster0.0vthl.mongodb.net/?retryWrites=true&w=majority",
		"demo",
	)
	coll := NewImageCollection(db, "image")

	id, _ := uuid.NewUUID()
	name := fmt.Sprintf("hello_%d.png", id.ID())
	info := entity.ImageInfo{
		UserName:  "user_abc",
		ImagePath: "/user_abc/" + name,
		FileName:  name,
	}

	err := coll.Save(info)
	assert.Nil(t, err)

	// re-check
	cursor, err := coll.client.Find(context.TODO(), bson.D{{"name", name}})
	require.NoError(t, err)
	requireCursorLength(t, cursor, 1)
}

func requireCursorLength(t *testing.T, cursor *mongo.Cursor, length int) {
	i := 0
	for cursor.Next(context.Background()) {
		i++
	}

	require.NoError(t, cursor.Err())
	require.Equal(t, i, length)
}
