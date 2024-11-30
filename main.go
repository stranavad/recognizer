package main

import (
	"recognizer/db"
	"recognizer/exam"
	"recognizer/files"
	"recognizer/game"
	"recognizer/group"
	"recognizer/item"
	"recognizer/types"
	"recognizer/user"

	_ "github.com/breml/rootcerts"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		userId, err := user.ParseToken(tokenString)
		if err != nil {
			c.JSON(403, "Invalid token")
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

func main() {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length", "Content-Type", "Authorization"},
	}))

	config := types.ServiceConfig{
		DB: db.GetDB(),
		S3: files.GetS3Client(),
	}

	/*
		User
	*/
	userService := user.NewUserService(config)
	userGroup := r.Group("/user")
	userGroup.GET("current", AuthMiddleware(), userService.GetCurrentUser)
	userGroup.POST("create", userService.CreateUser)
	userGroup.POST("login", userService.LoginUser)

	/*
		Exams
	*/
	examService := exam.NewExamService(config)
	examGroups := r.Group("/exam")
	examGroups.Use(AuthMiddleware())
	examGroups.POST("", examService.CreateExam)
	examGroups.PUT(":examId", examService.UpdateExam)
	examGroups.GET(":examId", examService.GetExam)
	examGroups.GET("/stats/:examId", examService.GetExamStats)
	examGroups.DELETE(":examId", examService.DeleteExam)
	examGroups.GET("", examService.ListExams)

	/*
		Groups
	*/
	groupService := group.NewGroupService(config)
	groupGroup := r.Group("/group")
	groupGroup.Use(AuthMiddleware())
	groupGroup.GET("by-exam/:examId", groupService.ListGroups)
	groupGroup.POST("", groupService.CreateGroup)
	groupGroup.PUT(":groupId", groupService.UpdateGroup)
	groupGroup.DELETE(":groupId", groupService.DeleteGroup)

	/*
		Items
	*/
	itemsService := item.NewItemService(config)
	itemGroup := r.Group("/items")
	itemGroup.Use(AuthMiddleware())
	itemGroup.POST("", itemsService.CreateItem)
	itemGroup.PUT(":itemId", itemsService.UpdateItem)
	itemGroup.GET(":itemId", itemsService.GetItem)
	itemGroup.DELETE(":itemId", itemsService.DeleteItem)
	itemGroup.GET("/by-exam/:examId", itemsService.ListItems)

	/*
		Game
	*/
	gameService := game.NewGameService(config)
	gameGroup := r.Group("/game")
	gameGroup.Use(AuthMiddleware())
	gameGroup.GET("/:examId", gameService.GetItem)
	gameGroup.POST("/result", gameService.GetResult)

	/*
		Files
	*/
	filesService := files.NewFilesService(config)

	r.POST("/file", AuthMiddleware(), filesService.UploadFile)

	err := r.Run()

	if err != nil {
		panic("Could not start gin server")
	}
}
