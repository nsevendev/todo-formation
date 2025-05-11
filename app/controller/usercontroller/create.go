package usercontroller

import (
	"fmt"
	"todof/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

// Create godoc
// @Summary Créer un utilisateur
// @Description Création d’un nouvel utilisateur
// @Tags user
// @Accept json
// @Produce json
// @Param user body auth.UserCreateDto true "DTO d'inscription utilisateur"
// @Success 200 {object} ginresponse.JsonFormatterSwag "Utilisateur créé avec succès"
// @Failure 400 {object} ginresponse.JsonFormatterSwag "Erreur de validation"
// @Failure 500 {object} ginresponse.JsonFormatterSwag "Erreur interne"
// @Router /user/register [post]
func (u *userController) Create(c *gin.Context) {
	var userCreateDto auth.UserCreateDto
	if err := c.ShouldBindJSON(&userCreateDto); err != nil {
		logger.Ef("Erreur de validation: %v", err)
		ginresponse.BadRequest(c, "Erreur de validation", ginresponse.ErrorModel{
			Message: err.Error(),
			Type:    "Validation",
			Detail:  fmt.Sprintf("%v", err),
		})
		return
	}

	_, err := u.userService.Register(c, userCreateDto)
	if err != nil {
		logger.Ef("%v", err)
		ginresponse.InternalServerError(c, err.Error(), err.Error())
		return
	}

	ginresponse.Created(c, "Utilisateur créé avec succès", []string{})
}
