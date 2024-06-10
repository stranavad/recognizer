package exam

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"recognizer/db"
	"recognizer/types"
	"strconv"
)

type Service struct {
	types.ServiceConfig
}

func NewExamService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) CreateExam(c *gin.Context) {
	var data types.CreateExamDto

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the exam already exists by name
	var foundExam *db.Exam
	service.DB.Where("name = ?", data.Name).First(&foundExam)

	if foundExam.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Exam with this name already exists"})
		return
	}

	// Create exam
	createdExam := db.Exam{
		Name: data.Name,
	}
	service.DB.Create(&createdExam)

	// Load fields from DB
	service.DB.First(&createdExam)

	c.JSON(200, createdExam)
}

func (service *Service) UpdateExam(c *gin.Context) {
	// Load exam by ID
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundExam *db.Exam
	res := service.DB.First(&foundExam, uint(examIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"message": `Exam not found`})
		return
	}

	// Get update request
	var data types.CreateExamDto

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundExam.Name = data.Name
	service.DB.Save(&foundExam)

	// Load fields from DB
	service.DB.First(&foundExam)

	c.JSON(200, foundExam)
}

func (service *Service) GetExam(c *gin.Context) {
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundExam *db.Exam
	res := service.DB.First(&foundExam, uint(examIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Exam not found"})
		return
	}

	c.JSON(200, foundExam)
}

func (service *Service) DeleteExam(c *gin.Context) {
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundExam *db.Exam
	res := service.DB.First(&foundExam, uint(examIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Exam not found"})
		return
	}

	service.DB.Delete(&foundExam)

	c.JSON(200, gin.H{"message": "Successfully deleted exam"})
}

func (service *Service) ListExams(c *gin.Context) {
	var exams []db.Exam

	service.DB.Find(&exams)

	c.JSON(200, exams)
}
