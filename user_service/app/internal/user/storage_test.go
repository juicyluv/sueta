package user_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/juicyluv/sueta/user_service/app/config"
	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/juicyluv/sueta/user_service/app/internal/user/apperror"
	"github.com/juicyluv/sueta/user_service/app/internal/user/db"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/juicyluv/sueta/user_service/app/pkg/mongo"
	"github.com/stretchr/testify/assert"
)

func NewTestStorage(t *testing.T) (user.Storage, func() error) {
	logger.Init()

	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "user_service") {
		wd = filepath.Dir(wd)
	}

	cfg := config.Get(fmt.Sprintf("%s/app/config/test.yml", wd), fmt.Sprintf("%s/.env", wd))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongoClient, err := mongo.NewMongoClient(ctx, cfg.DB.Database, cfg.DB.URL)
	if err != nil {
		t.Fatalf("cannot connect to mongodb: %v", err)
	}

	userStorage := db.NewStorage(mongoClient, cfg.DB.Collection)

	teardown := func() error {
		return mongoClient.Collection(cfg.DB.Collection).Drop(context.Background())
	}

	return userStorage, teardown
}

func TestUserStorage_Create(t *testing.T) {
	storage, teardown := NewTestStorage(t)
	defer func() { assert.NoError(t, teardown()) }()

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

func TestUserStorage_FindByEmail(t *testing.T) {
	storage, teardown := NewTestStorage(t)
	defer func() { assert.NoError(t, teardown()) }()

	users := []user.User{
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

	for _, u := range users {
		createStorageUser(t, storage, &u)
	}

	testCases := []struct {
		name          string
		email         string
		expectedError error
	}{
		{
			name:          "existing user 1",
			email:         users[0].Email,
			expectedError: nil,
		},
		{
			name:          "existing user 2",
			email:         users[1].Email,
			expectedError: nil,
		},
		{
			name:          "existing user 3",
			email:         users[2].Email,
			expectedError: nil,
		},
		{
			name:          "user does not exist",
			email:         "someemail@example.com",
			expectedError: apperror.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := storage.FindByEmail(context.Background(), tc.email)
			assert.EqualValues(t, err, tc.expectedError)
			if tc.expectedError == nil {
				found := false
				for _, v := range users {
					if v.Email == u.Email {
						found = true
					}
				}
				assert.True(t, found)
			}
		})
	}
}

func TestUserStorage_FindById(t *testing.T) {
	storage, teardown := NewTestStorage(t)
	defer func() { assert.NoError(t, teardown()) }()

	users := []user.User{
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

	for i := range users {
		id := createStorageUser(t, storage, &users[i])
		users[i].UUID = id
	}

	testCases := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "existing user 1",
			id:            users[0].UUID,
			expectedError: nil,
		},
		{
			name:          "existing user 2",
			id:            users[1].UUID,
			expectedError: nil,
		},
		{
			name:          "existing user 3",
			id:            users[2].UUID,
			expectedError: nil,
		},
		{
			name:          "user does not exist",
			id:            "6202fbbf5ff5d5a12a72a195",
			expectedError: apperror.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := storage.FindById(context.Background(), tc.id)
			assert.EqualValues(t, err, tc.expectedError)
			if tc.expectedError == nil {
				found := false
				for _, v := range users {
					if v.UUID == u.UUID {
						found = true
					}
				}
				assert.True(t, found)
			}
		})
	}
}

func TestUserStorage_UpdatePartially(t *testing.T) {
	storage, teardown := NewTestStorage(t)
	defer func() { assert.NoError(t, teardown()) }()

	u := user.User{
		Email:    "test1@mail.com",
		Username: "test1",
		Password: "qwerty",
	}

	u.UUID = createStorageUser(t, storage, &u)

	testCases := []struct {
		name          string
		expectedError error
		beforeTest    func()
		afterTest     func(t *testing.T)
	}{
		{
			name:          "email update",
			expectedError: nil,
			beforeTest: func() {
				u.Email = "updated@example.com"
			},
			afterTest: func(t *testing.T) {
				found, err := storage.FindByEmail(context.Background(), u.Email)
				assert.NoError(t, err)
				assert.Equal(t, found.UUID, u.UUID)
				assert.Equal(t, found.Email, "updated@example.com")
			},
		},
		{
			name:          "username update",
			expectedError: nil,
			beforeTest: func() {
				u.Username = "test123"
			},
			afterTest: func(t *testing.T) {
				found, err := storage.FindByEmail(context.Background(), u.Email)
				assert.NoError(t, err)
				assert.Equal(t, found.UUID, u.UUID)
				assert.Equal(t, found.Username, "test123")
			},
		},
		{
			name:          "password update",
			expectedError: nil,
			beforeTest: func() {
				u.Password = "juicyluv"
				u.HashPassword()
			},
			afterTest: func(t *testing.T) {
				found, err := storage.FindByEmail(context.Background(), u.Email)
				assert.NoError(t, err)
				assert.Equal(t, found.UUID, u.UUID)
				assert.True(t, u.ComparePassword("juicyluv"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.beforeTest()
			err := storage.UpdatePartially(context.Background(), &u)
			assert.EqualValues(t, err, tc.expectedError)
			tc.afterTest(t)
		})
	}
}

func createStorageUser(t *testing.T, storage user.Storage, user *user.User) string {
	id, err := storage.Create(context.Background(), user)
	assert.NoError(t, err)
	return id
}
