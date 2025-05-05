package termsconditionscontroller

import (
	"todof/internal/termsconditions"

	"github.com/gin-gonic/gin"
)

type termsConditionsController struct {
	termsConditionsService termsconditions.TermsConditionsServiceInterface
}

type TermsConditionsControllerInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
}

func NewTermsConditionsController(termsConditionsService termsconditions.TermsConditionsServiceInterface) TermsConditionsControllerInterface {
	return &termsConditionsController{
		termsConditionsService: termsConditionsService,
	}
}



