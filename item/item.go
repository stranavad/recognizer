package item

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

func NewItemService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) CreateItem(c *gin.Context) {
	var data types.CreateItem

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemToCreate := db.Item{
		Name:    data.Name,
		Image:   data.Image,
		GroupId: data.GroupId,
		ExamId:  data.ExamId,
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
	res := service.DB.First(&foundItem, uint(itemIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	foundItem.Name = data.Name
	foundItem.GroupId = data.GroupId
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
	res := service.DB.First(&foundItem, uint(itemIdParam))

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
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
