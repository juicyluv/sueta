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
