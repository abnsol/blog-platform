package routers

import (
	"github.com/blog-platform/delivery/controllers"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/repositories"
	"github.com/blog-platform/usecases"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(repositories.DB)
	ei := infrastructure.NewSMTPEmailService()
	pi := infrastructure.NewPasswordInfrastructure()
	uu := usecases.NewUserUsecase(ur, ei, pi)
	uc := controllers.NewUserController(uu)

	group.POST("/register", uc.Register)
}