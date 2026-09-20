package http

import (
	"net/http"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/service/users"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	serv users.UserSv
}

func NewUserController(serv users.UserSv) *UserController {
	return &UserController{serv: serv}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := c.serv.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user_id": id})
}

func (c *UserController) VerifyUserEmail(ctx *gin.Context) {
	var verify models.VerifyUserEmail
	if err := ctx.ShouldBindJSON(&verify); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.serv.VerfiyUserEmail(ctx, &verify); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to verify user", "error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) Login(ctx *gin.Context) {
	var user models.UserLogin
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userData, token, err := c.serv.Login(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  userData,
		"token": token,
	})
}
