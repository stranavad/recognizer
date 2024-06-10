package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"recognizer/db"
	"recognizer/exam"
	"recognizer/group"
	"recognizer/types"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length", "Content-Type", "Authorization"},
	}))

	config := types.ServiceConfig{
		DB: db.GetDB(),
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

	err := r.Run()

	if err != nil {
		panic("Could not start gin server")
	}
}
