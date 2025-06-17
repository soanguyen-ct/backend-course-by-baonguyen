package mongostore

import (
	"ct-backend-course-baonguyen/internal/entity"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCollection_ChangePassword(t *testing.T) {
	db := MustDatabase(
		"mongodb+srv://bao:baopass@cluster0.0vthl.mongodb.net/?retryWrites=true&w=majority",
		"demo",
	)
	userColl := NewUserCollection(db, "user")

	err := userColl.Create(entity.UserInfo{
		Username: fmt.Sprintf("name_%s.png", uuid.New().String()),
		Password: "abc_123",
	})
	assert.Nil(t, err)

	// TODO #4 Write test case ChangePassword

	t.Run("Duplicate", func(t *testing.T) {

	})

	t.Run("Success", func(t *testing.T) {

	})
}
