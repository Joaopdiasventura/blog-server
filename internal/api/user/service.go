package user

import (
	"context"
	"net/http"

	"github.com/joaopdias/blog-server/internal/shared/auth"
	"github.com/joaopdias/blog-server/internal/shared/errors"
)

type UserService struct {
	repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repository: repo}
}

func (s *UserService) Create(ctx context.Context, createUserDTO CreateUserDTO) (User, string, *errors.ApiError) {
	_, err := s.repository.FindByEmail(ctx, createUserDTO.Email)
	if err == nil {
		return User{}, "", errors.NewApiError(http.StatusConflict, "user already exists")
	}

	createUserDTO.Password = auth.HashPassword(createUserDTO.Password)

	user, err := s.repository.Create(ctx, createUserDTO)
	if err != nil {
		return User{}, "", errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	token, err := auth.GenerateJWT(user.Id)
	if err != nil {
		return User{}, "", errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	user.Password = ""

	return user, token, nil
}

func (s *UserService) Login(ctx context.Context, loginUserDTO LoginUserDTO) (User, string, *errors.ApiError) {
	user, err := s.repository.FindByEmail(ctx, loginUserDTO.Email)
	if err != nil {
		return User{}, "", errors.NewApiError(http.StatusNotFound, "user not found")
	}

	if !auth.CheckPasswordHash(loginUserDTO.Password, user.Password) {
		return User{}, "", errors.NewApiError(http.StatusUnauthorized, "wrong password")
	}

	token, err := auth.GenerateJWT(user.Id)
	if err != nil {
		return User{}, "", errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	user.Password = ""

	return user, token, nil
}

func (s *UserService) FindById(ctx context.Context, id string) (User, *errors.ApiError) {
	user, err := s.repository.FindById(ctx, id)
	if err != nil {
		return User{}, errors.NewApiError(http.StatusNotFound, "user not found")
	}
	return user, nil
}

func (s *UserService) DecodeToken(ctx context.Context, token string) (User, *errors.ApiError) {
	id, err := auth.ParseJWT(token)
	if err != nil {
		return User{}, errors.NewApiError(http.StatusUnauthorized, "invalid token")
	}

	return s.FindById(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id string, updateUserDTO UpdateUserDTO) (User, *errors.ApiError) {
	u, apiErr := s.FindById(ctx, id)
	if apiErr != nil {
		return u, apiErr
	}

	if updateUserDTO.Password != nil && *updateUserDTO.Password != "" {
		hashed := auth.HashPassword(*updateUserDTO.Password)
		updateUserDTO.Password = &hashed
	}

	user, err := s.repository.Update(ctx, id, updateUserDTO)
	if err != nil {
		return User{}, errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) *errors.ApiError {
	_, apiErr := s.FindById(ctx, id)
	if apiErr != nil {
		return apiErr
	}

	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

