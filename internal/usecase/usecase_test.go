package usecase

import (
	"bytes"
	"context"
	"ct-backend-course-baonguyen/internal/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestUcImplement_UploadImage(t *testing.T) {
	err := errors.New("fail connection")
	testCases := []struct {
		name         string
		imageBucket  ImageBucket
		expImagePath string
		expError     error
	}{
		{
			name:         "success",
			imageBucket:  &mockImgBucket{name: "/static/images/imageTest.png"},
			expImagePath: "/static/images/imageTest.png",
			expError:     nil,
		},
		{
			name:         "failed image bucket",
			imageBucket:  &mockImgBucket{Err: err},
			expImagePath: "",
			expError:     err,
		},
	}

	for _, tc := range testCases {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUseCase(nil, nil, tt.imageBucket)
			image := bytes.NewReader([]byte("Hello"))
			imageInfo, err := uc.UploadImage(context.TODO(), &entity.UploadImageRequest{Content: image})
			if tt.expError != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, imageInfo.ImageUrl, tt.expImagePath)

		})
	}

}

type mockUserStore struct {
}

func (m *mockUserStore) Save(info entity.UserInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockUserStore) Get(username string) (entity.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

type mockImgBucket struct {
	Err  error
	name string
}

func (m *mockImgBucket) SaveImage(ctx context.Context, name string, r io.Reader) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}

	return m.name, nil
}
