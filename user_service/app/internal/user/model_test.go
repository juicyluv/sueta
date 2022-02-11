package user_test

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserDTO_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		input         *user.CreateUserDTO
		expectedError *validation.Errors
	}{
		{
			name: "valid input",
			input: &user.CreateUserDTO{
				Email:          "hello@mail.ru",
				Username:       "qwerty",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedError: nil,
		},
		{
			name: "invalid email",
			input: &user.CreateUserDTO{
				Email:          "hellomail.ru",
				Username:       "qwerty",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{"email": errors.New("must be a valid email address")},
		},
		{
			name: "username lenght is less than 3",
			input: &user.CreateUserDTO{
				Email:          "hello@mail.ru",
				Username:       "qw",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{"username": errors.New("the length must be between 3 and 20")},
		},
		{
			name: "username lenght is greater than 20",
			input: &user.CreateUserDTO{
				Email:          "hello@mail.ru",
				Username:       "verylongusernameeeeee",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{"username": errors.New("the length must be between 3 and 20")},
		},
		{
			name: "invalid email and username",
			input: &user.CreateUserDTO{
				Email:          "hellomail.ru",
				Username:       "verylongusernameeeeee",
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{
				"email":    errors.New("must be a valid email address"),
				"username": errors.New("the length must be between 3 and 20"),
			},
		},
		{
			name: "username and email are not provided",
			input: &user.CreateUserDTO{
				Password:       "qwerty",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{
				"email":    errors.New("cannot be blank"),
				"username": errors.New("cannot be blank"),
			},
		},
		{
			name: "password length is less than 6",
			input: &user.CreateUserDTO{
				Email:          "hello@mail.ru",
				Username:       "username",
				Password:       "qwe",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{
				"password": errors.New("the length must be between 6 and 24"),
			},
		},
		{
			name: "password length is greater than 24",
			input: &user.CreateUserDTO{
				Email:          "hello@mail.ru",
				Username:       "username",
				Password:       "qwertyqwertyqwertyqwertywwww",
				RepeatPassword: "qwerty",
			},
			expectedError: &validation.Errors{
				"password": errors.New("the length must be between 6 and 24"),
			},
		},
		{
			name: "repeatPassword is not provided",
			input: &user.CreateUserDTO{
				Email:    "hello@mail.ru",
				Username: "username",
				Password: "qwerty",
			},
			expectedError: &validation.Errors{
				"repeatPassword": errors.New("cannot be blank"),
			},
		},
		{
			name:  "empty input",
			input: &user.CreateUserDTO{},
			expectedError: &validation.Errors{
				"email":          errors.New("cannot be blank"),
				"username":       errors.New("cannot be blank"),
				"password":       errors.New("cannot be blank"),
				"repeatPassword": errors.New("cannot be blank"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.expectedError != nil {
				assert.EqualValues(t, err, *tc.expectedError)
			}
		})
	}
}

func TestUser_HashPassword(t *testing.T) {
	testCases := []struct {
		name string
		u    user.User
	}{
		{
			name: "valid password",
			u:    user.User{Password: "qwerty"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
