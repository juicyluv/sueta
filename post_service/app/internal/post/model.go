package post

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Post struct {
	UUID      string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserUUID  string    `json:"userId"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	Comments  []Comment `json:"comments"`
}

type CreatePostDTO struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	UserUUID string `json:"userId"`
}

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (p *CreatePostDTO) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(
			&p.Title,
			is.ASCII,
			validation.Length(3, 200),
			validation.Required,
		),
		validation.Field(
			&p.Content,
			is.ASCII,
			validation.Length(10, 5000),
			validation.Required,
		),
		validation.Field(
			&p.UserUUID,
			validation.Required,
		),
	)
}

type UpdatePostDTO struct {
	UUID     string  `json:"id"`
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	UserUUID *string `json:"userId"`
}

// Validate will validates current struct fields.
// Returns an error if something doesn't fit rules.
func (p *UpdatePostDTO) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(
			&p.Title,
			is.ASCII,
			validation.Length(3, 200),
		),
		validation.Field(
			&p.Content,
			is.ASCII,
			validation.Length(10, 5000),
		),
	)
}

type Comment struct {
	UUID      string `json:"id"`
	Content   string `json:"content"`
	UserUUID  string `json:"userId"`
	Verified  bool   `json:"verified"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
