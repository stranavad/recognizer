package group

import (
	"errors"
	"net/http"
	"recognizer/db"
	"recognizer/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	types.ServiceConfig
}

func NewGroupService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) CreateGroup(c *gin.Context) {
	var data types.CreateGroupDto

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check exam existence
	var foundExam *db.Exam
	res := service.DB.Preload("Groups").First(&foundExam, data.ExamID)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Exam not found"})
		return
	}

	// Authorization check
	if foundExam.UserID != c.MustGet("userId").(uint){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if group with this name already exists
	groupExists := false
	for _, value := range foundExam.Groups {
		if value.Name == data.Name {
			groupExists = true
			break
		}
	}

	if groupExists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group with this name already exists in this exams"})
		return
	}

	createdGroup := db.Group{
		Name:   data.Name,
		ExamID: data.ExamID,
	}

	// Create and load group
	service.DB.Create(&createdGroup)
	service.DB.First(&createdGroup)

	c.JSON(200, createdGroup)
}

func (service *Service) UpdateGroup(c *gin.Context) {
	groupIdParam, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var foundGroup *db.Group
	res := service.DB.Preload("Exam").First(&foundGroup, uint(groupIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
		return
	}

	if foundGroup.Exam.UserID != c.MustGet("userId").(uint){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind request body
	var data types.UpdateGroupDto
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	foundGroup.Name = data.Name
	service.DB.Save(&foundGroup)
	service.DB.First(&foundGroup)

	c.JSON(200, foundGroup)
}

func (service *Service) ListGroups(c *gin.Context) {
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var groups []db.Group
	service.DB.Where("exam_id = ?", uint(examIdParam)).Find(&groups)

	c.JSON(200, groups)
}

func (service *Service) DeleteGroup(c *gin.Context) {
	groupIdParam, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var foundGroup *db.Group
	res := service.DB.Preload("Exam").First(&foundGroup, uint(groupIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
		return
	}

	if foundGroup.Exam.UserID != c.MustGet("userId").(uint){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	service.DB.Delete(&foundGroup)
	service.DB.Where("group_id = ?", foundGroup.ID).Delete(&db.Item{})
}
