package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaopdias/blog-server/internal/api/post"
	"github.com/joaopdias/blog-server/internal/api/user"
)

type services struct {
	userService    *user.UserService
	postController *post.PostController
	userController *user.UserController
}

func NewRouter(pool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	register(r, wire(pool))
	return r
}

func wire(pool *pgxpool.Pool) *services {
	userRepository := user.NewPostgresUserRepository(pool)
	userService := user.NewUserService(userRepository)
	userController := user.NewUserController(userService)

	postRepo := post.NewPostgresPostRepository(pool)
	postService := post.NewPostService(postRepo, userService)
	postController := post.NewPostController(postService)

	return &services{
		userService:    userService,
		postController: postController,
		userController: userController,
	}
}

func register(r *gin.Engine, s *services) {
	s.userController.RegisterRoutes(r)
	s.postController.RegisterRoutes(r)
}
