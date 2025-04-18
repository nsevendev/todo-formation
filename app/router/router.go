package router

import (
	"todof/app/controller/taskcontroller"
	"todof/internal/models"
	"todof/internal/taskservice"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func Router(r *gin.Engine) {
	modelManager := models.NewManager()

	taskService := taskservice.NewTaskService(modelManager)
	
	taskController := taskcontroller.NewTaskController(taskService)

	v1 := r.Group("api/v1")

	// task
	v1Task := v1.Group("/task")
	v1Task.GET("/:id", taskController.GetOneById)

	r.NoRoute(func(ctx *gin.Context) {
		logger.Wf("Route inconnue : %s %s", ctx.Request.Method, ctx.Request.URL.Path)
		ginresponse.NotFound(ctx, "La route demandée n'existe pas.", "La route demandée n'existe pas.")
		ctx.Abort()
	})
}