package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *PostService
}

func NewPostController(service *PostService) *PostController {
	return &PostController{service: service}
}

func (c *PostController) RegisterRoutes(r *gin.Engine) {
	r.POST("/post", c.Create)
	r.GET("/post", c.FindById)
	r.GET("/post/findMany", c.FindMany)
	r.GET("/post/findAllByAuthor", c.FindAllByAuthor)
	r.DELETE("/post", c.Delete)
}

func (c *PostController) Create(ctx *gin.Context) {
	var dto CreatePostDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "invalid body"})
		return
	}

	post, err := c.service.Create(ctx.Request.Context(), dto)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "post created",
		"post":    post,
	})
}

func (c *PostController) FindById(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing id"})
		return
	}

	post, err := c.service.FindById(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (c *PostController) FindMany(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid offset"})
		return
	}

	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, apiErr := c.service.FindMany(ctx.Request.Context(), limit, offset)
	if apiErr != nil {
		ctx.AbortWithStatusJSON(apiErr.Code, gin.H{"message": apiErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "posts found",
		"posts":   posts,
	})
}

func (c *PostController) FindAllByAuthor(ctx *gin.Context) {
	author := ctx.Query("author")
	if author == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing author"})
		return
	}

	posts, apiErr := c.service.FindAllByAuthor(ctx.Request.Context(), author)
	if apiErr != nil {
		ctx.AbortWithStatusJSON(apiErr.Code, gin.H{"message": apiErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "posts found",
		"posts":   posts,
	})
}

func (c *PostController) Delete(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing id"})
		return
	}

	err := c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "post deleted",
	})
}
