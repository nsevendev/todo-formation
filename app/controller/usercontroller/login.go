package usercontroller

import (
	"fmt"
	"todof/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

// Login godoc
// @Summary Authentifier un utilisateur
// @Description Authentification d’un utilisateur
// @Tags user
// @Accept json
// @Produce json
// @Param user body auth.UserLoginDto true "DTO d'authentification utilisateur"
// @Success 200 {object} ginresponse.JsonFormatterSwag "Connexion réussie"
// @Failure 400 {object} ginresponse.JsonFormatterSwag "Erreur de validation"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur d'authentification"
// @Router /user/login [post]
func (u *userController) Login(c *gin.Context) {
	var userLoginDto auth.UserLoginDto
	if err := c.ShouldBindJSON(&userLoginDto); err != nil {
		logger.Ef("Erreur de validation: %v", err)
		ginresponse.BadRequest(c, "Erreur de validation", ginresponse.ErrorModel{
			Message: err.Error(),
			Type:    "Validation",
			Detail:  fmt.Sprintf("%v", err),
		})
		return
	}

	token, err := u.userService.Login(c, userLoginDto)
	if err != nil {
		logger.Ef("%v", err)
		ginresponse.InternalServerError(c, err.Error(), err.Error())
		return
	}

	ginresponse.Success(c, "Connexion réussie", map[string]string{
		"token": token,
	})
}
