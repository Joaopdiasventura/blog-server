package post

import (
	"context"
	"net/http"

	"github.com/joaopdias/blog-server/internal/api/user"
	"github.com/joaopdias/blog-server/internal/shared/errors"
)

type PostService struct {
	repository  PostRepository
	userService *user.UserService
}

func NewPostService(repo PostRepository, userService *user.UserService) *PostService {
	return &PostService{repository: repo, userService: userService}
}

func (s *PostService) Create(ctx context.Context, createPostDTO CreatePostDTO) (Post, *errors.ApiError) {
	_, apiErr := s.userService.FindById(ctx, createPostDTO.AuthorId)
	if apiErr != nil {
		return Post{}, errors.NewApiError(http.StatusBadRequest, "author does not exist")
	}

	post, err := s.repository.Create(ctx, createPostDTO)
	if err != nil {
		return Post{}, errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	return post, nil
}

func (s *PostService) FindById(ctx context.Context, id string) (Post, *errors.ApiError) {
	post, err := s.repository.FindById(ctx, id)
	if err != nil {
		return Post{}, errors.NewApiError(http.StatusNotFound, "post not found")
	}

	return post, nil
}

func (s *PostService) FindMany(ctx context.Context, limit, offset int) ([]Post, *errors.ApiError) {
	posts, err := s.repository.FindMany(ctx, limit, offset)
	if err != nil {
		return nil, errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	return posts, nil
}

func (s *PostService) FindAllByAuthor(ctx context.Context, author string) ([]Post, *errors.ApiError) {
	posts, err := s.repository.FindAllByAuthor(ctx, author)
	if err != nil {
		return nil, errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	return posts, nil
}

func (s *PostService) Delete(ctx context.Context, id string) *errors.ApiError {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
