package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	userURL  = "/api/users/:uuid"
	usersURL = "/api/users"
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

	h, ok := handler.(*user.Handler)
	if !ok {
		t.Fatal("cannot convert handler to user.Handler type")
	}

	type input struct {
		Email          string `json:"email,omitempty"`
		Username       string `json:"username,omitempty"`
		Password       string `json:"password,omitempty"`
		RepeatPassword string `json:"repeatPassword,omitempty"`
		ExtraField     string `json:"extra,omitempty"`
	}

	type createResponse struct {
		Id string `json:"id"`
	}

	testCases := []struct {
		name                  string
		expectedCode          int
		input                 input
		expectedResponse      createResponse
		expectedErrorResponse *apperror.AppError
	}{
		{
			name:         "valid input",
			expectedCode: http.StatusCreated,
			input: input{
				Email:          "test@mail.com",
				Username:       "test",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: nil,
		},
		{
			name:         "extra field input",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "test",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
				ExtraField:     "something",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"request body contains unknown key \"extra\"",
				"invalid request body",
			),
		},
		{
			name:         "empty email",
			expectedCode: http.StatusBadRequest,
			input: input{
				Username:       "test",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"email: cannot be blank.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "empty email and password",
			expectedCode: http.StatusBadRequest,
			input: input{
				Username:       "test",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"email: cannot be blank; password: cannot be blank.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "passwords dont match",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "test",
				Password:       "qwerty",
				RepeatPassword: "qwerty123",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"passwords don't match",
				"provided passwords must to match",
			),
		},
		{
			name:         "invalid email",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "testmail.com",
				Username:       "test",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"email: must be a valid email address.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "username less than 3",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "te",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"username: the length must be between 3 and 20.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "username greater than 20",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "qwerqwerqwerqwerqwerqwerqwer",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"username: the length must be between 3 and 20.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "username is not alphanumeric",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "antosha_44",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"username: must contain English letters and digits only.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "password is not alphanumeric",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "test",
				Password:       "aNTOsha_44",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"password: must contain English letters and digits only.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "password length less than 6",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "test",
				Password:       "qwe",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"password: the length must be between 6 and 24.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "password length greater than 24",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "test",
				Password:       "qwertyqwertyqwertyqwertywwww",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				"password: the length must be between 6 and 24.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "empty input",
			expectedCode: http.StatusBadRequest,
			input:        input{},
			expectedErrorResponse: apperror.BadRequestError(
				"email: cannot be blank; password: cannot be blank; repeatPassword: cannot be blank; username: cannot be blank.",
				"input validation failed. please, provide valid values",
			),
		},
		{
			name:         "email taken",
			expectedCode: http.StatusBadRequest,
			input: input{
				Email:          "test@mail.com",
				Username:       "andrew",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedErrorResponse: apperror.BadRequestError(
				apperror.ErrEmailTaken.Error(),
				"",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(&tc.input)
			if err != nil {
				t.Fatalf("cannot marshal input: %v", err)
			}

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, usersURL, bytes.NewBuffer(body))
			assert.NoError(t, err)

			h.CreateUser(rec, req)
			res := rec.Result()

			assert.Equal(t, tc.expectedCode, res.StatusCode)

			response, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.expectedErrorResponse != nil {
				expectedResponse, err := json.Marshal(tc.expectedErrorResponse)
				assert.NoError(t, err)

				assert.NoError(t, err)
				assert.EqualValues(t, response, expectedResponse)
			} else {
				var createdId createResponse
				err := json.Unmarshal(response, &createdId)
				assert.NoError(t, err)
				assert.NotEmpty(t, createdId.Id)
			}
		})
	}
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

func TestUserHandler_GetUserByEmailAndPassword(t *testing.T) {
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
		expectedCode          int
		expectedResponse      *user.User
		expectedErrorResponse *apperror.AppError
		email                 string
		password              string
	}{
		{
			name:         "user exists",
			expectedCode: http.StatusOK,
			expectedResponse: &user.User{
				UUID:         id,
				Email:        u.Email,
				Username:     u.Username,
				Verified:     false,
				RegisteredAt: time.Now().UTC().Format("2006/01/02"),
			},
			expectedErrorResponse: nil,
			email:                 u.Email,
			password:              u.Password,
		},
		{
			name:             "user not found by email",
			expectedCode:     http.StatusNotFound,
			expectedResponse: nil,
			expectedErrorResponse: apperror.NewAppError(
				http.StatusNotFound,
				"requested resource is not found",
				"please, double check your request",
			),
			email:    "test2@mail.com",
			password: u.Password,
		},
		{
			name:                  "user not found by wrong password",
			expectedCode:          http.StatusBadRequest,
			expectedResponse:      nil,
			expectedErrorResponse: apperror.BadRequestError(apperror.ErrWrongPassword.Error(), ""),
			email:                 u.Email,
			password:              "qwerty123",
		},
		{
			name:             "empty email",
			expectedCode:     http.StatusBadRequest,
			expectedResponse: nil,
			expectedErrorResponse: apperror.BadRequestError(
				"empty email or password",
				"email and password must be provided",
			),
			email:    "",
			password: u.Password,
		},
		{
			name:             "empty password",
			expectedCode:     http.StatusBadRequest,
			expectedResponse: nil,
			expectedErrorResponse: apperror.BadRequestError(
				"empty email or password",
				"email and password must be provided",
			),
			email:    u.Email,
			password: "",
		},
		{
			name:             "empty email and password",
			expectedCode:     http.StatusBadRequest,
			expectedResponse: nil,
			expectedErrorResponse: apperror.BadRequestError(
				"empty email or password",
				"email and password must be provided",
			),
			email:    "",
			password: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s?email=%s&password=%s", usersURL, tc.email, tc.password),
				nil,
			)
			assert.NoError(t, err)

			h.GetUserByEmailAndPassword(rec, req)
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

func TestUserHandler_UpdatePartially(t *testing.T) {
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

	stringPtr := func(s string) *string {
		return &s
	}

	testCases := []struct {
		name                  string
		expectedCode          int
		input                 *user.UpdateUserDTO
		expectedErrorResponse *apperror.AppError
	}{
		{
			name:         "valid email update",
			expectedCode: http.StatusOK,
			input: &user.UpdateUserDTO{
				Email:       stringPtr("test@mail.com"),
				OldPassword: &u.Password,
			},
			expectedErrorResponse: nil,
		},
		{
			name:         "invalid email update",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				Email:       stringPtr("testmail.com"),
				OldPassword: &u.Password,
			},
			expectedErrorResponse: apperror.BadRequestError(
				"email: must be a valid email address.",
				"you have provided invalid values",
			),
		},
		{
			name:         "valid username update",
			expectedCode: http.StatusOK,
			input: &user.UpdateUserDTO{
				Username:    stringPtr("test2"),
				OldPassword: &u.Password,
			},
			expectedErrorResponse: nil,
		},
		{
			name:         "username less than 3 update",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				Username:    stringPtr("qw"),
				OldPassword: &u.Password,
			},
			expectedErrorResponse: apperror.BadRequestError(
				"username: the length must be between 3 and 20.",
				"you have provided invalid values",
			),
		},
		{
			name:         "username greater than 20 update",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				Username:    stringPtr("qwertyqwertyqwerty123123"),
				OldPassword: &u.Password,
			},
			expectedErrorResponse: apperror.BadRequestError(
				"username: the length must be between 3 and 20.",
				"you have provided invalid values",
			),
		},
		{
			name:         "valid new password update",
			expectedCode: http.StatusOK,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("qwertyqwerty"),
				OldPassword: &u.Password,
			},
			expectedErrorResponse: nil,
		},
		{
			name:         "new password less than 6 update",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("qwert"),
				OldPassword: stringPtr("qwertyqwerty"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				"newPassword: the length must be between 6 and 24.",
				"you have provided invalid values",
			),
		},
		{
			name:         "new password greater than 24 update",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("qwertyqwertyqwertyqwerty123123"),
				OldPassword: stringPtr("qwertyqwerty"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				"newPassword: the length must be between 6 and 24.",
				"you have provided invalid values",
			),
		},
		{
			name:         "invalid old password",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("qwerty"),
				OldPassword: stringPtr("qwertyqwerty123"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				apperror.ErrWrongPassword.Error(),
				"you entered wrong password",
			),
		},
		{
			name:         "old password is not provided",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("qwerty"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				"oldPassword: cannot be blank.",
				"you have provided invalid values",
			),
		},
		{
			name:         "old password is not provided",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("qwerty"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				"oldPassword: cannot be blank.",
				"you have provided invalid values",
			),
		},
		{
			name:         "username is not alphanumeric",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				Username:    stringPtr("antoha_44ru"),
				OldPassword: stringPtr("qwertyqwerty"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				"username: must contain English letters and digits only.",
				"you have provided invalid values",
			),
		},
		{
			name:         "new password is not alphanumeric",
			expectedCode: http.StatusBadRequest,
			input: &user.UpdateUserDTO{
				NewPassword: stringPtr("antoha_44ru"),
				OldPassword: stringPtr("qwertyqwerty"),
			},
			expectedErrorResponse: apperror.BadRequestError(
				"newPassword: must contain English letters and digits only.",
				"you have provided invalid values",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("cannot marshal input: %v", err)
			}

			rec := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPatch, userURL, bytes.NewBuffer(body))
			assert.NoError(t, err)
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
				{Key: "uuid", Value: id},
			})
			req = req.WithContext(ctx)

			h.UpdateUserPartially(rec, req)
			res := rec.Result()

			assert.Equal(t, res.StatusCode, tc.expectedCode)

			response, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			if tc.expectedErrorResponse != nil {
				expectedResponse, err := json.Marshal(tc.expectedErrorResponse)
				assert.NoError(t, err)

				assert.NoError(t, err)
				assert.EqualValues(t, response, expectedResponse)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
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

	u1 := user.CreateUserDTO{
		Email:          "test@mail.com",
		Username:       "test",
		Password:       "qwerty",
		RepeatPassword: "qwerty",
	}

	id1, err := createUser(h, &u1)
	assert.NoError(t, err)

	u2 := user.CreateUserDTO{
		Email:          "test2@mail.com",
		Username:       "test2",
		Password:       "qwerty",
		RepeatPassword: "qwerty",
	}

	id2, err := createUser(h, &u2)
	assert.NoError(t, err)

	testCases := []struct {
		name                  string
		expectedCode          int
		id                    string
		expectedErrorResponse *apperror.AppError
	}{
		{
			name:                  "delete user 1",
			expectedCode:          http.StatusOK,
			id:                    id1,
			expectedErrorResponse: nil,
		},
		{
			name:         "delete user 1 again",
			expectedCode: http.StatusNotFound,
			id:           id1,
			expectedErrorResponse: apperror.NewAppError(
				http.StatusNotFound,
				"requested resource is not found",
				"please, double check your request",
			),
		},
		{
			name:                  "delete user 2",
			expectedCode:          http.StatusOK,
			id:                    id2,
			expectedErrorResponse: nil,
		},
		{
			name:         "delete user 2 again",
			expectedCode: http.StatusNotFound,
			id:           id1,
			expectedErrorResponse: apperror.NewAppError(
				http.StatusNotFound,
				"requested resource is not found",
				"please, double check your request",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, userURL, nil)
			assert.NoError(t, err)
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
				{Key: "uuid", Value: tc.id},
			})
			req = req.WithContext(ctx)

			h.DeleteUser(rec, req)
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
				rec := httptest.NewRecorder()

				req, err := http.NewRequest(http.MethodDelete, userURL, nil)
				assert.NoError(t, err)
				ctx := req.Context()
				ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
					{Key: "uuid", Value: tc.id},
				})
				req = req.WithContext(ctx)

				h.GetUser(rec, req)
				res := rec.Result()

				assert.Equal(t, res.StatusCode, http.StatusNotFound)
			}
		})
	}
}

func createUser(h *user.Handler, u *user.CreateUserDTO) (string, error) {

	body, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, usersURL, bytes.NewBuffer(body))
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
