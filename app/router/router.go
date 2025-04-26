package router

import (
	"os"
	"todof/app/controller/taskcontroller"
	"todof/app/controller/usercontroller"
	"todof/internal/auth"
	"todof/internal/task"

	initializer "todof/internal/init"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
)

func Router(r *gin.Engine) {
	userRepo := auth.NewUserRepo(initializer.Db)
	userService := auth.NewUserService(userRepo, os.Getenv("JWT_SECRET"))
	authMiddle := auth.NewAuthMiddleware(userService)
	userController := usercontroller.NewUserController(userService)

	taskRepo := task.NewTaskRepo(initializer.Db)
	taskService := task.NewTaskService(taskRepo)
	taskController := taskcontroller.NewTaskController(taskService, userService)

	v1 := r.Group("api/v1")

	v1Task := v1.Group("/task")
	v1Task.Use(authMiddle.RequireAuth())
	v1Task.POST("/", taskController.Create)
	v1Task.GET("/", taskController.GetAllByUser)
	v1Task.PUT("/:id/done/user", taskController.UpdateOneDonePropertyByUser)
	v1Task.PUT("/:id/user", taskController.UpdateOneLabelPropertyByUser)
	v1Task.DELETE("/:id/user", taskController.DeleteOneByUser)
	v1Task.POST("/:id/user", taskController.DeleteOneByUser)
	v1Task.POST("/delete/user", taskController.DeleteManyByUser)

	v1User := v1.Group("/user")
	v1User.POST("/register", userController.Create)
	v1User.POST("/login", userController.Login)
	v1User.GET("/profil", authMiddle.RequireAuth(), userController.GetProfilCurrentUser)

	r.NoRoute(func(ctx *gin.Context) {
		logger.Wf("Route inconnue : %s %s", ctx.Request.Method, ctx.Request.URL.Path)
		ginresponse.NotFound(ctx, "La route demandée n'existe pas.", "La route demandée n'existe pas.")
		ctx.Abort()
	})
}