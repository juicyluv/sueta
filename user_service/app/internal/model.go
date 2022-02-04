package internal

type User struct {
	UUID         string `json:"uuid" bson:"_id,omitempty"`
	Username     string `json:"username", bson:"username,omitempty"`
	Email        string `json:"email" bson:"email,omitempty"`
	Password     string `json:"-" bson:"password,omitempty"`
	Verified     bool   `json:"verified" bson:"verified,omitempty"`
	RegisteredAt string `json:"registeredAt" bson:"registeredAt,omitempty"`
	Role         Role   `json:"role"`
}

type Role struct {
	UUID string `json:"uuid" bson:"_id,omitempty"`
	Role string `json:"role" bson:"role,omitempty"`
}
