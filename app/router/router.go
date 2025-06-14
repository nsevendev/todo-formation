package router

import (
	"github.com/nsevenpack/env/env"
	"todof/app/controller/taskcontroller"
	"todof/app/controller/usercontroller"
	"todof/docs"
	"todof/internal/auth"
	"todof/internal/task"

	initializer "todof/internal/init"

	"github.com/gin-gonic/gin"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(r *gin.Engine) {
	pathApiV1 := "api/v1"

	userRepo := auth.NewUserRepo(initializer.Db)
	userService := auth.NewUserService(userRepo, env.Get("JWT_SECRET"))
	authMiddle := auth.NewAuthMiddleware(userService)
	userController := usercontroller.NewUserController(userService)

	taskRepo := task.NewTaskRepo(initializer.Db)
	taskService := task.NewTaskService(taskRepo, userRepo)
	taskController := taskcontroller.NewTaskController(taskService, userService)

	docs.SwaggerInfo.BasePath = "/" + pathApiV1
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group(pathApiV1)

	v1Task := v1.Group("/task")
	v1Task.Use(authMiddle.RequireAuth())
	v1Task.POST("/", taskController.Create)
	v1Task.GET("/", taskController.GetAllByUser)
	v1Task.PUT("/:id/done/user", taskController.UpdateOneDonePropertyByUser)
	v1Task.PUT("/:id/label/user", taskController.UpdateOneLabelPropertyByUser)
	v1Task.DELETE("/:id/user", taskController.DeleteOneByUser)
	v1Task.POST("/delete/user", taskController.DeleteManyByUser)
	v1Task.POST("/delete/tasks", authMiddle.RequireAuth(), authMiddle.RequireRole("admin"), taskController.DeleteById)
	v1Task.DELETE("/delete/all", authMiddle.RequireAuth(), authMiddle.RequireRole("admin"), taskController.DeleteAllTasks)

	v1User := v1.Group("/user")
	v1User.POST("/register", userController.Create)
	v1User.POST("/login", userController.Login)
	v1User.GET("/profil", authMiddle.RequireAuth(), userController.GetProfilCurrentUser)
	v1User.DELETE("/", authMiddle.RequireAuth(), userController.DeleteOneByUser)
	v1User.DELETE("/users", authMiddle.RequireAuth(), authMiddle.RequireRole("admin"), userController.DeleteByAdmin)
	v1User.DELETE("/users/all", authMiddle.RequireAuth(), authMiddle.RequireRole("admin"), userController.DeleteAllByAdmin)

	r.NoRoute(func(ctx *gin.Context) {
		logger.Wf("Route inconnue : %s %s", ctx.Request.Method, ctx.Request.URL.Path)
		ginresponse.NotFound(ctx, "La route demandée n'existe pas.", "La route demandée n'existe pas.")
		ctx.Abort()
	})
}
