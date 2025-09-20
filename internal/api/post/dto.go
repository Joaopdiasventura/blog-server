package post

type CreatePostDTO struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	AuthorId    string `json:"authorId" binding:"required"`
}
