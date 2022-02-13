package user_test

import (
	"context"
	"errors"
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

func TestUserService_CreateUser(t *testing.T) {
	logger.Init()
	l := logger.GetLogger()

	userStorage, teardown := NewTestStorage(t)
	service := user.NewService(userStorage, l)

	testCases := []struct {
		name          string
		expectedError error
		input         *user.CreateUserDTO
		testBehaviour func(id string, err error)
	}{
		{
			name:          "valid input",
			expectedError: nil,
			input: &user.CreateUserDTO{
				Email:          "test@mail.ru",
				Username:       "test",
				Password:       "passwod",
				RepeatPassword: "password",
			},
			testBehaviour: func(id string, err error) {
				assert.Greater(t, len(id), 0)
				assert.NoError(t, err)
			},
		},
		{
			name:          "email taken",
			expectedError: apperror.ErrEmailTaken,
			input: &user.CreateUserDTO{
				Email:          "test@mail.ru",
				Username:       "test2",
				Password:       "passwod",
				RepeatPassword: "password",
			},
			testBehaviour: func(id string, err error) {
				assert.Equal(t, len(id), 0)
				assert.Error(t, err)
				assert.ErrorIs(t, err, apperror.ErrEmailTaken)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Nil(t, tc.input.Validate())

			id, err := service.Create(context.Background(), tc.input)
			tc.testBehaviour(id, err)
		})
	}

	assert.NoError(t, teardown())
}

func TestUserService_GetByEmailAndPassword(t *testing.T) {
	logger.Init()
	l := logger.GetLogger()

	userStorage, teardown := NewTestStorage(t)
	service := user.NewService(userStorage, l)

	testCases := []struct {
		name          string
		expectedError error
		email         string
		password      string
	}{
		{
			name:          "valid input",
			expectedError: nil,
			email:         "test@mail.com",
			password:      "qwerty",
		},
		{
			name:          "invalid email",
			expectedError: apperror.ErrNoRows,
			email:         "test2@mail.com",
			password:      "qwerty",
		},
		{
			name:          "invalid password",
			expectedError: apperror.ErrWrongPassword,
			email:         "test@mail.com",
			password:      "qwerty2",
		},
		{
			name:          "invalid email and password",
			expectedError: apperror.ErrNoRows,
			email:         "test2@mail.com",
			password:      "qwerty2",
		},
	}

	created := &user.CreateUserDTO{
		Email:    "test@mail.com",
		Username: "test",
		Password: "qwerty",
	}

	id1, err := service.Create(context.Background(), created)
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := service.GetByEmailAndPassword(context.Background(), tc.email, tc.password)
			assert.EqualValues(t, err, tc.expectedError)
			if tc.expectedError == nil {
				assert.Equal(t, id1, u.UUID)
				assert.Equal(t, created.Email, u.Email)
				assert.Equal(t, created.Username, u.Username)
				assert.True(t, u.ComparePassword(created.Password))
			}
		})
	}

	assert.NoError(t, teardown())
}

func TestUserService_GetById(t *testing.T) {
	logger.Init()
	l := logger.GetLogger()

	userStorage, teardown := NewTestStorage(t)
	service := user.NewService(userStorage, l)

	created := &user.CreateUserDTO{
		Email:          "test@mail.com",
		Username:       "test",
		Password:       "qwerty",
		RepeatPassword: "qwerty",
	}

	id, err := service.Create(context.Background(), created)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		uuid          string
		expectedError error
	}{
		{
			name:          "valid id",
			uuid:          id,
			expectedError: nil,
		},
		{
			name:          "no user with given id",
			uuid:          id[:len(id)-1] + "7",
			expectedError: apperror.ErrNoRows,
		},
		{
			name:          "invalid uuid",
			uuid:          "invaliduuid",
			expectedError: errors.New("failed to find user by uuid: failed to convert hex to objectid: the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := service.GetById(context.Background(), tc.uuid)
			assert.EqualValues(t, err, tc.expectedError)
			if tc.expectedError == nil {
				assert.Equal(t, u.UUID, id)
			}
		})
	}

	assert.NoError(t, teardown())
}

func TestUserService_UpdatePartially(t *testing.T) {
	logger.Init()
	l := logger.GetLogger()

	userStorage, teardown := NewTestStorage(t)
	service := user.NewService(userStorage, l)

	created := &user.CreateUserDTO{
		Email:          "test@mail.com",
		Username:       "test",
		Password:       "qwerty",
		RepeatPassword: "qwerty",
	}

	stringPtr := func(s string) *string {
		return &s
	}

	id, err := service.Create(context.Background(), created)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		input         *user.UpdateUserDTO
		expectedError error
	}{
		{
			name: "update email",
			input: &user.UpdateUserDTO{
				UUID:        id,
				Email:       stringPtr("newemail@gmail.com"),
				OldPassword: stringPtr("qwerty"),
			},
			expectedError: nil,
		},
		{
			name: "update username",
			input: &user.UpdateUserDTO{
				UUID:        id,
				Username:    stringPtr("newusername"),
				OldPassword: stringPtr("qwerty"),
			},
			expectedError: nil,
		},
		{
			name: "update password",
			input: &user.UpdateUserDTO{
				UUID:        id,
				NewPassword: stringPtr("qwerty123123"),
				OldPassword: stringPtr("qwerty"),
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualValues(t, tc.expectedError, service.UpdatePartially(context.Background(), tc.input))
		})
	}

	u, err := service.GetById(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, "newemail@gmail.com")
	assert.Equal(t, u.Username, "newusername")
	assert.True(t, u.ComparePassword("qwerty123123"))

	assert.NoError(t, teardown())
}
