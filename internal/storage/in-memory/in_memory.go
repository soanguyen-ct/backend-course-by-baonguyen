package inmemory

import (
	"ct-backend-course-baonguyen/internal/entity"
	"errors"
	"sync"
)

func NewUserStore() *userStore {
	return &userStore{data: make(map[string]entity.UserInfo)}
}

type userStore struct {
	mu   sync.Mutex
	data map[string]entity.UserInfo
}

func (u *userStore) Save(info entity.UserInfo) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Save the user information in the data map
	u.data[info.Username] = info

	return nil
}

func (u *userStore) Get(username string) (entity.UserInfo, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Retrieve the user information from the data map
	user, found := u.data[username]
	if !found {
		return entity.UserInfo{}, ErrUserNotFound
	}

	return user, nil
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserExisted = errors.New("user existed")
