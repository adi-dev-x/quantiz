package handlers

type AddBlog struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	Title        string  `json:"title" validate:"required"`
	Descriptions string  `json:"descriptions" validate:"required"`
	Author       string  `json:"author" validate:"required"`
	Content      string  `json:"content" validate:"required"`
	Categories   []int64 `json:"categories" validate:"required"`
}
type UpdateBlogRequest struct {
	Title        *string `json:"title"`
	Descriptions *string `json:"descriptions"`
	Author       *string `json:"author"`
	Content      *string `json:"content"`
	Categories   []int64 `json:"categories"`
}
