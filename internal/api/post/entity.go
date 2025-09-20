package post

import (
	"time"

	"github.com/joaopdias/blog-server/internal/api/user"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorId  string    `json:"authorId"`
	Author    user.User `json:"author,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
