package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) RegisterRoutes(r *gin.Engine) {
	r.POST("/user", c.Create)
	r.POST("/user/login", c.Login)
	r.GET("/user/decodeToken", c.DecodeToken)
	r.PATCH("/user", c.Update)
	r.DELETE("/user", c.Delete)
}

func (c *UserController) Create(ctx *gin.Context) {
	var dto CreateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	user, token, err := c.service.Create(ctx.Request.Context(), dto)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
		"token":   token,
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	var dto LoginUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	user, token, err := c.service.Login(ctx.Request.Context(), dto)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user logged in",
		"user":    user,
		"token":   token,
	})
}

func (c *UserController) DecodeToken(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing token"})
		return
	}

	user, err := c.service.DecodeToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "missing id"})
		return
	}

	var dto UpdateUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	user, err := c.service.Update(ctx.Request.Context(), id, dto)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user updated",
		"user":    user,
	})
}

func (c *UserController) Delete(ctx *gin.Context) {
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
		"message": "user deleted",
	})
}
