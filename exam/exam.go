package exam

import (
	"errors"
	"fmt"
	"math"
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

func NewExamService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) CreateExam(c *gin.Context) {
	var data types.CreateExamDto

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(uint)

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
		UserID: userId,
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

	userId := c.MustGet("userId").(uint)

	var foundExam *db.Exam
	res := service.DB.First(&foundExam, uint(examIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"message": `Exam not found`})
		return
	}

	if foundExam.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": `Unathorized`})
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

type LeaderboardItem struct {
	UserID    uint   `json:"userId"`
	Correct   int    `json:"correct"`
	Wrong     int    `json:"wrong"`
	Nickname  string `json:"nickname"`
	Total     int    `json:"total"`
	Points    int    `json:"points"`
	Percentage int    `json:"percentage"`
}

func(service  *Service) GetExamStats(c *gin.Context){
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data []LeaderboardItem
	err = service.DB.Model(&db.ScorePoint{}).
		Select("score_points.user_id, "+
			"CAST(SUM(CASE WHEN score_points.correct THEN 1 ELSE 0 END) AS INT) as correct, "+
			"CAST(SUM(CASE WHEN score_points.correct THEN 0 ELSE 1 END) AS INT) as wrong, "+
			"users.username",
		).
		Where("score_points.exam_id = ?", uint(examIdParam)).
		Joins("INNER JOIN users ON score_points.user_id = users.id").
		Group("score_points.user_id, users.username").
		Scan(&data).Error

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying database"})
		return
	}

	for i, item := range data {
		total := item.Correct + item.Wrong
		data[i].Total = total
		data[i].Points = (item.Correct * 10) - (item.Wrong * 5) + 100
		data[i].Percentage = int(math.Round((float64(item.Correct) / float64(total)) * 100))
	}

	// Sort by points in descending order (using a simple bubble sort for this example)
	for i := 0; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Points < data[j].Points {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

}

func (service *Service) DeleteExam(c *gin.Context) {
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(uint)

	var foundExam *db.Exam
	res := service.DB.First(&foundExam, uint(examIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Exam not found"})
		return
	}

	if foundExam.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
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
