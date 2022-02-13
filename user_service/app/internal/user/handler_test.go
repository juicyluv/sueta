package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/juicyluv/sueta/user_service/app/internal/handler"
	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/juicyluv/sueta/user_service/app/internal/user/apperror"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const (
	userURL = "/api/users/:uuid"
)

func NewTestHandler(t *testing.T) (handler.Handling, func() error) {
	logger.Init()
	l := logger.GetLogger()

	router := httprouter.New()

	userStorage, teardown := NewTestStorage(t)
	service := user.NewService(userStorage, l)
	handler := user.NewHandler(l, service)
	handler.Register(router)

	return handler, teardown
}

func TestUserHandler_CreateUser(t *testing.T) {
	handler, teardown := NewTestHandler(t)
	defer func() {
		assert.NoError(t, teardown())
	}()
	router := httprouter.New()
	handler.Register(router)

}

func TestUserHandler_GetUser(t *testing.T) {
	handler, teardown := NewTestHandler(t)
	defer func() {
		assert.NoError(t, teardown())
	}()
	router := httprouter.New()
	handler.Register(router)

	h, ok := handler.(*user.Handler)
	if !ok {
		t.Fatal("cannot convert handler to user.Handler type")
	}

	u := user.CreateUserDTO{
		Email:          "test@mail.com",
		Username:       "test",
		Password:       "qwerty",
		RepeatPassword: "qwerty",
	}

	id, err := createUser(h, &u)
	assert.NoError(t, err)

	testCases := []struct {
		name                  string
		uuid                  string
		expectedCode          int
		expectedErrorResponse *apperror.AppError
		expectedResponse      user.User
	}{
		{
			name:                  "found",
			uuid:                  id,
			expectedCode:          http.StatusOK,
			expectedErrorResponse: nil,
			expectedResponse: user.User{
				UUID:         id,
				Email:        u.Email,
				Username:     u.Username,
				Verified:     false,
				RegisteredAt: time.Now().UTC().Format("2006/01/02"),
			},
		},
		{
			name:                  "not found",
			uuid:                  "62056f8cf21b83383a5ae7fa",
			expectedCode:          http.StatusNotFound,
			expectedErrorResponse: apperror.ErrNotFound,
		},
		{
			name:                  "invalid uuid",
			uuid:                  "invaliduuid",
			expectedCode:          http.StatusBadRequest,
			expectedErrorResponse: apperror.NewAppError(400, apperror.ErrInvalidUUID.Error(), ""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, userURL, nil)
			assert.NoError(t, err)
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
				{Key: "uuid", Value: tc.uuid},
			})
			req = req.WithContext(ctx)

			h.GetUser(rec, req)
			res := rec.Result()

			assert.Equal(t, res.StatusCode, tc.expectedCode)

			response, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			if tc.expectedErrorResponse != nil {
				expectedResponse, err := json.Marshal(tc.expectedErrorResponse)
				assert.NoError(t, err)

				assert.NoError(t, err)
				assert.EqualValues(t, response, expectedResponse)
			} else {
				expectedResponseBytes, err := json.Marshal(&tc.expectedResponse)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBytes, response)
			}
		})
	}
}

func createUser(h *user.Handler, u *user.CreateUserDTO) (string, error) {

	body, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, userURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	rec := httptest.NewRecorder()

	h.CreateUser(rec, req)

	type response struct {
		Id string `json:"id"`
	}

	res := response{}

	err = json.NewDecoder(rec.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}
