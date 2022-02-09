package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model.
type User struct {
	UUID         string `json:"uuid" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email,omitempty"`
	Username     string `json:"username" bson:"username,omitempty"`
	Password     string `json:"-" bson:"password,omitempty"`
	Verified     bool   `json:"verified" bson:"verified,omitempty"`
	RegisteredAt string `json:"registeredAt" bson:"registeredAt,omitempty"`
	Role         Role   `json:"role"`
}

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
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (u *CreateUserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, is.Email, validation.Required),
		validation.Field(&u.Username, is.Alphanumeric, validation.Required),
		validation.Field(&u.Password, is.Alphanumeric, validation.Required),
	)
}

// UpdateUserDTO is used to update user record.
type UpdateUserDTO struct {
	UUID        string
	Email       *string `json:"email"`
	Username    *string `json:"username"`
	OldPassword *string `json:"oldPassword"`
	NewPassword *string `json:"newPassword"`
	RoleUUID    *string `json:"roleId"`
}

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (u *UpdateUserDTO) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, is.Email),
		validation.Field(&u.Username, is.Alphanumeric),
		validation.Field(&u.OldPassword, is.Alphanumeric),
		validation.Field(&u.NewPassword, is.Alphanumeric),
		validation.Field(&u.RoleUUID, is.UUIDv4),
	)
}

// Role represents user role model.
type Role struct {
	UUID string `json:"uuid" bson:"_id,omitempty"`
	Role string `json:"role" bson:"role,omitempty"`
}
