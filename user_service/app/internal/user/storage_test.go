package user_test

import (
	"context"
	"testing"

	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserStorage_NewStorage(t *testing.T) {
	storage, _ := NewTestStorage(t)

	var _ user.Storage = storage
	assert.NotNil(t, storage)
}

func TestUserStorage_Create(t *testing.T) {
	storage, teardown := NewTestStorage(t)
	defer teardown()

	testCases := []user.User{
		{
			Email:    "test1@mail.com",
			Username: "test1",
			Password: "qwerty",
		},
		{
			Email:    "test2@mail.com",
			Username: "test2",
			Password: "qwerty123",
		},
		{
			Email:    "test3@mail.com",
			Username: "test3",
			Password: "qwertyyyy",
		},
	}

	for _, tc := range testCases {
		rawPassword := tc.Password

		err := tc.HashPassword()
		assert.NoError(t, err)

		assert.True(t, tc.ComparePassword(rawPassword))

		id, err := storage.Create(context.Background(), &tc)
		assert.NoError(t, err)
		assert.True(t, len(id) > 0)

		u, err := storage.FindById(context.Background(), id)
		assert.NotNil(t, u)
		assert.NoError(t, err)
		assert.Equal(t, tc.Email, u.Email)
		assert.Equal(t, tc.Username, u.Username)
	}
}
