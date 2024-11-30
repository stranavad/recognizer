package game

import (
	"math/rand"
	"net/http"
	"recognizer/db"
	"recognizer/types"
	"strconv"

	"github.com/gin-gonic/gin"
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

	var similarAnswers []string
	for itemIndex, item := range items {
		if itemIndex != randomIndex && item.GroupID == randomItem.GroupID {
			similarAnswers = append(similarAnswers, item.Name)
		}
	}

	// Now we'll shuffle these elements
	rand.Shuffle(len(similarAnswers), func(i, j int) { similarAnswers[i], similarAnswers[j] = similarAnswers[j], similarAnswers[i] })
	itemsToGet := len(similarAnswers)
	if itemsToGet > 3 {
		itemsToGet = 3
	}

	answersItems := similarAnswers[:itemsToGet]
	answersItems = append(answersItems, randomItem.Name)

	// Shuffle the answers
	rand.Shuffle(len(answersItems), func(i, j int) { answersItems[i], answersItems[j] = answersItems[j], answersItems[i] })

	response := types.GameResponse{
		ItemId:  randomItem.ID,
		Image:   randomItem.Image,
		Answers: answersItems,
	}

	c.JSON(http.StatusOK, response)
}

func (service *Service) GetResult(c *gin.Context) {
	var data types.GetResult

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item *db.Item
	service.DB.First(&item, data.ItemId)

	isCorrect := item.Name == data.Answer

	scorePoint := db.ScorePoint{
		UserID: c.MustGet("userId").(uint),
		ExamID: item.ExamID,
		ItemID: item.ID,
		Correct: isCorrect,
	}
	service.DB.Create(scorePoint)

	c.JSON(http.StatusOK, gin.H{"correct": isCorrect})
}
