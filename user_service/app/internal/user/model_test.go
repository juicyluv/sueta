package user_test

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserDTO_Validate(t *testing.T) {
	t.Parallel()

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

func TestUpdateUserDTO_Validate(t *testing.T) {
	t.Parallel()

	type in struct {
		Email       string
		Username    string
		OldPassword string
		NewPassword string
	}

	testCases := []struct {
		name          string
		input         in
		expectedError *validation.Errors
	}{
		{
			name: "valid input",
			input: in{
				Email:       "test@mail.com",
				Username:    "test",
				OldPassword: "qwerty",
				NewPassword: "qwerty123",
			},
			expectedError: nil,
		},
		{
			name: "invalid email",
			input: in{
				Email:       "testmail.com",
				Username:    "test",
				OldPassword: "qwerty",
				NewPassword: "qwerty123",
			},
			expectedError: &validation.Errors{"email": errors.New("must be a valid email address")},
		},
		{
			name: "invalid username",
			input: in{
				Email:       "test@mail.com",
				Username:    "qwer_ty",
				OldPassword: "qwerty",
				NewPassword: "qwerty123",
			},
			expectedError: &validation.Errors{"username": errors.New("must contain English letters and digits only")},
		},
		{
			name: "new password less than 6",
			input: in{
				Email:       "test@mail.com",
				Username:    "test",
				OldPassword: "qwerty",
				NewPassword: "qwe",
			},
			expectedError: &validation.Errors{"newPassword": errors.New("the length must be between 6 and 24")},
		},
		{
			name: "new password greater than 24",
			input: in{
				Email:       "test@mail.com",
				Username:    "test",
				OldPassword: "qwerty",
				NewPassword: "qwertyqwertyqwertyqwertywwww",
			},
			expectedError: &validation.Errors{"newPassword": errors.New("the length must be between 6 and 24")},
		},
		{
			name: "empty string input",
			input: in{
				Email:       "",
				Username:    "",
				OldPassword: "",
				NewPassword: "",
			},
			expectedError: &validation.Errors{"oldPassword": errors.New("cannot be blank")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			update := &user.UpdateUserDTO{
				Email:       &tc.input.Email,
				Username:    &tc.input.Username,
				OldPassword: &tc.input.OldPassword,
				NewPassword: &tc.input.NewPassword,
			}

			err := update.Validate()
			if tc.expectedError != nil {
				assert.EqualValues(t, err, *tc.expectedError)
			}
		})
	}

	t.Run("empty nil string", func(t *testing.T) {
		update := &user.UpdateUserDTO{
			Email:       nil,
			Username:    nil,
			OldPassword: nil,
			NewPassword: nil,
		}

		expectedError := validation.Errors{"oldPassword": errors.New("cannot be blank")}

		err := update.Validate()
		assert.EqualValues(t, err, expectedError)
	})
}

func TestUser_HashPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		u    user.User
	}{
		{
			name: "valid password",
			u:    user.User{Password: "qwerty"},
		},
		{
			name: "long password",
			u:    user.User{Password: "asdfasdfasdfa123"},
		},
		{
			name: "empty password",
			u:    user.User{Password: ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.u.HashPassword()
			assert.NoError(t, err)
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		u    user.User
	}{
		{
			name: "valid password",
			u:    user.User{Password: "qwerty"},
		},
		{
			name: "long password",
			u:    user.User{Password: "asdfasdfasdfa123"},
		},
		{
			name: "empty password",
			u:    user.User{Password: ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			password := tc.u.Password
			err := tc.u.HashPassword()
			assert.NoError(t, err)
			assert.True(t, tc.u.ComparePassword(password))
		})
	}
}
