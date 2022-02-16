package post

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

type UpdatePostDTO struct {
	UUID     string  `json:"id"`
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	UserUUID *string `json:"userId"`
}

type Comment struct {
	UUID      string `json:"id"`
	Content   string `json:"content"`
	UserUUID  string `json:"userId"`
	Verified  bool   `json:"verified"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
