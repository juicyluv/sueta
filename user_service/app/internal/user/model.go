package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model.
type User struct {
	UUID         string `json:"uuid" bson:"_id,omitempty" example:"1"`
	Email        string `json:"email" bson:"email,omitempty" example:"1"`
	Username     string `json:"username" bson:"username,omitempty" example:"1"`
	Password     string `json:"-" bson:"password,omitempty" example:"1"`
	Verified     bool   `json:"verified" bson:"verified,omitempty" example:"1"`
	RegisteredAt string `json:"registeredAt" bson:"registeredAt,omitempty" example:"1"`
} // @name User

// HashPassword will encrypt current user password.
// Returns an error on failure.
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ComparePassword compares hashed user password with given raw password.
// If it doesn't match, returns false.
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// CreateUserDTO is used to create user.
type CreateUserDTO struct {
	Email          string `json:"email"`
	Username       string `json:"username" minLength:"3" maxLength:"20"`
	Password       string `json:"password" minLength:"6" maxLength:"24"`
	RepeatPassword string `json:"repeatPassword"`
} // @name CreateUserInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (u *CreateUserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Email,
			is.Email,
			validation.Required,
		),
		validation.Field(
			&u.Username,
			is.Alphanumeric,
			validation.Length(3, 20),
			validation.Required,
		),
		validation.Field(
			&u.Password,
			is.Alphanumeric,
			validation.Length(6, 24),
			validation.Required,
		),
	)
}

// UpdateUserDTO is used to update user record.
type UpdateUserDTO struct {
	UUID        string  `json:"-"`
	Email       *string `json:"email"`
	Username    *string `json:"username" minLength:"3" maxLength:"20"`
	OldPassword *string `json:"oldPassword"`
	NewPassword *string `json:"newPassword" minLength:"6" maxLength:"24"`
} // @name UpdateUserInput

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (u *UpdateUserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, is.Email),
		validation.Field(&u.Username, is.Alphanumeric),
		validation.Field(&u.OldPassword, is.Alphanumeric, validation.Required),
		validation.Field(
			&u.NewPassword,
			validation.Length(6, 24),
			is.Alphanumeric),
	)
}
