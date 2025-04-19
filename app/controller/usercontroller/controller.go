package usercontroller

import (
	"todof/internal/auth"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService auth.UserServiceInterface
}

type UserControllerInterface interface {
	Create(c *gin.Context)
	Login(c *gin.Context)
	GetProfilCurrentUser(c *gin.Context)
}

func NewUserController(userService auth.UserServiceInterface) UserControllerInterface {
	return &userController{
		userService: userService,
	}
}
