package http

import (
	"net/http"
	"os"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/service/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type UserController struct {
	serv users.UserSv
}

func NewUserController(serv users.UserSv) *UserController {
	return &UserController{serv: serv}
}

func (c *UserController) GoogleAuth(ctx *gin.Context) {
	var googleAuthReq models.GoogleAuthReq

	if err := ctx.ShouldBindJSON(&googleAuthReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	payload, err := idtoken.Validate(ctx, googleAuthReq.IDToken, os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate Google ID token"})
		return
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	googleID, _ := payload.Claims["sub"].(string)
	picture, _ := payload.Claims["picture"].(string)

	user := &models.User{
		Email:          email,
		Name:           name,
		GoogleID:       googleID,
		ProfilePicture: picture,
	}

	userData, token, err := c.serv.AuthUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  userData,
		"token": token,
	})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	user_id, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid access token",
			"message": err.Error(),
		})
		return
	}

	userID := uuid.MustParse(user_id)

	user, err := c.serv.GetUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get user",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) Refresh(ctx *gin.Context) {
	user_id, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid access token",
			"message": err.Error(),
		})
		return
	}

	userID := uuid.MustParse(user_id)

	refreshToken, err := GetRefreshToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to fetch refresh token",
			"message": err.Error(),
		})
		return
	}

	accessTokens, err := c.serv.RefreshToken(ctx, userID, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to create new access tokens",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, accessTokens)
}
