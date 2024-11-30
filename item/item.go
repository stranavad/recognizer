package item

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

func NewItemService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) CreateItem(c *gin.Context) {
	var data types.CreateItem


	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First try to find the exam
	var exam *db.Exam
	res := service.DB.Preload("Groups").First(&exam, data.ExamId)

	if errors.Is(res.Error, gorm.ErrRecordNotFound){
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	if exam.UserID != c.MustGet("userId").(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Then we check the group
	found := false
    for _, item := range exam.Groups {
        if item.ID == data.GroupId {
            found = true
            break
        }
    }

    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
        return
    }

	itemToCreate := db.Item{
		Name:    data.Name,
		Image:   data.Image,
		GroupID: data.GroupId,
		ExamID:  data.ExamId,
	}

	service.DB.Create(&itemToCreate)
	service.DB.First(&itemToCreate)

	c.JSON(200, itemToCreate)
}

func (service *Service) UpdateItem(c *gin.Context) {
	itemIdParam, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data types.UpdateItem

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundItem *db.Item
	res := service.DB.Preload("Exam").First(&foundItem, uint(itemIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Authorization check
	if foundItem.Exam.UserID != c.MustGet("userId").(uint){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	foundItem.Name = data.Name
	foundItem.GroupID = data.GroupId
	foundItem.Image = data.Image

	service.DB.Save(&foundItem)
	service.DB.First(&foundItem)

	c.JSON(200, foundItem)
}

func (service *Service) GetItem(c *gin.Context) {
	itemIdParam, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundItem *db.Item
	res := service.DB.First(&foundItem, uint(itemIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(200, foundItem)
}

func (service *Service) DeleteItem(c *gin.Context) {
	itemIdParam, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundItem *db.Item
	res := service.DB.Preload("Exam").First(&foundItem, uint(itemIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Authorization check
	if foundItem.Exam.UserID != c.MustGet("userId").(uint){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	service.DB.Delete(&foundItem)
	c.JSON(200, gin.H{"message": "Item deleted"})

}

func (service *Service) ListItems(c *gin.Context) {
	examIdParam, err := strconv.ParseUint(c.Param("examId"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []db.Item
	service.DB.Where("exam_id = ?", uint(examIdParam)).Find(&items)

	c.JSON(200, items)
}
