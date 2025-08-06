package controllers

import (
	"net/http"

	"github.com/blog-platform/domain"
	"github.com/gin-gonic/gin"
)

type UserRegisterDTO struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Identifier string `json:"identifier"`
	Password string `json:"password"`
}

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(uu domain.IUserUsecase) *UserController {
	return &UserController{
		userUsecase: uu,
	}
}

func (uc *UserController) Register(ctx *gin.Context)  {
	var userInput UserRegisterDTO

	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user := domain.User{
		Email: userInput.Email,
		Username: userInput.Username,
		Password: userInput.Password,
	}

	user, err := uc.userUsecase.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (uc *UserController) ActivateAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	err := uc.userUsecase.ActivateAccount(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "user account activated"})
}

func (uc *UserController) Login(ctx *gin.Context) {
	var userInput UserLoginDTO

	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	accessToken, refreshToken, err := uc.userUsecase.Login(userInput.Identifier, userInput.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access": accessToken,
		"refresh": refreshToken,
		"message": "Logged in successfully",
	})
}