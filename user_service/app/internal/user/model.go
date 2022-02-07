package user

type User struct {
	UUID         string `json:"uuid" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email,omitempty"`
	Username     string `json:"username" bson:"username,omitempty"`
	Password     string `json:"-" bson:"password,omitempty"`
	Verified     bool   `json:"verified" bson:"verified,omitempty"`
	RegisteredAt string `json:"registeredAt" bson:"registeredAt,omitempty"`
	Role         Role   `json:"role"`
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Email       *string `json:"email"`
	Username    *string `json:"username"`
	OldPassword *string `json:"oldPassword"`
	NewPassword *string `json:"newPassword"`
	RoleUUID    *string `json:"roleId"`
}

type Role struct {
	UUID string `json:"uuid" bson:"_id,omitempty"`
	Role string `json:"role" bson:"role,omitempty"`
}
