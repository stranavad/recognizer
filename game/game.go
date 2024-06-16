package game

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"recognizer/db"
	"recognizer/types"
	"strconv"
)

type Service struct {
	types.ServiceConfig
}

func NewGameService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) GetItem(c *gin.Context) {
	examIdParam, err := strconv.ParseInt(c.Param("examId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []*db.Item
	service.DB.Where("exam_id = ?", uint(examIdParam)).Find(&items)

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No items in this exam"})
		return
	}

	randomIndex := rand.Intn(len(items))
	randomItem := items[randomIndex]

	var similarItems []*db.Item
	for _, item := range items {
		if item.ID != randomItem.ID && item.GroupId == randomItem.GroupId {
			similarItems = append(similarItems, item)
		}
	}

	// Now we'll shuffle these elements
	rand.Shuffle(len(similarItems), func(i, j int) { similarItems[i], similarItems[j] = similarItems[i], similarItems[i] })
	itemsToGet := len(similarItems)
	if itemsToGet > 3 {
		itemsToGet = 3
	}

	answersItems := similarItems[:itemsToGet]

	answers := []string{randomItem.Name}
	for _, item := range answersItems {
		answers = append(answers, item.Name)
	}

	// Shuffle the answers
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[i], answers[i] })

	response := types.GameResponse{
		ItemId:  randomItem.ID,
		Image:   randomItem.Image,
		Answers: answers,
	}

	c.JSON(http.StatusOK, response)
}

func (service *Service) GetResult(c *gin.Context) {
	var data types.GetResult

	if err := c.ShouldBindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item *db.Item
	service.DB.First(&item, data.ItemId)

	c.JSON(http.StatusOK, gin.H{"correct": item.Name == data.Answer})
}
