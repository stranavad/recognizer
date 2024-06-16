package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"recognizer/db"
	"recognizer/exam"
	"recognizer/files"
	"recognizer/game"
	"recognizer/group"
	"recognizer/item"
	"recognizer/types"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

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
		Exams
	*/
	examService := exam.NewExamService(config)
	examGroups := r.Group("/exam")
	examGroups.POST("", examService.CreateExam)
	examGroups.PUT(":examId", examService.UpdateExam)
	examGroups.GET(":examId", examService.GetExam)
	examGroups.DELETE(":examId", examService.DeleteExam)
	examGroups.GET("", examService.ListExams)

	/*
		Groups
	*/
	groupService := group.NewGroupService(config)
	groupGroup := r.Group("/group")
	groupGroup.GET("by-exam/:examId", groupService.ListGroups)
	groupGroup.POST("", groupService.CreateGroup)
	groupGroup.PUT(":groupId", groupService.UpdateGroup)
	groupGroup.DELETE(":groupId", groupService.DeleteGroup)

	/*
		Items
	*/

	itemsService := item.NewItemService(config)
	itemGroup := r.Group("/items")
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
	gameGroup.GET("/:examId", gameService.GetItem)
	gameGroup.POST("/result", gameService.GetResult)

	/*
		Files
	*/
	filesService := files.NewFilesService(config)
	r.POST("/file", filesService.UploadFile)

	err := r.Run()

	if err != nil {
		panic("Could not start gin server")
	}
}
