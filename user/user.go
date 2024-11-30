package user

import (
	"fmt"
	"net/http"
	"recognizer/db"
	"recognizer/types"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type Service struct {
	types.ServiceConfig
}

func NewUserService(config types.ServiceConfig) Service {
	return Service{config}
}

var secret = "supersafesecret"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(userId uint) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 30 * 365).Unix() // Token will expire in a year
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(tokenString string) (uint, error) {
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}
	userId := uint((*claims)["user_id"].(float64))

	return userId, nil
}

func(service *Service) GetCurrentUser(c *gin.Context){
	userId := c.MustGet("userId").(uint)
	var foundUser db.User
	service.DB.First(&foundUser, userId)
	c.JSON(200, foundUser.ToSimpleUser())
}

func(service *Service) LoginUser(c *gin.Context){
	var request types.CreateUserDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser db.User
	err := service.DB.Where("username = ?", request.Username).First(&foundUser).Error
	if err != nil {
		c.JSON(403, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	if !checkPassword(request.Password, foundUser.Password) {
		c.JSON(403, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	token, tokenErr := CreateToken(foundUser.ID)

	if tokenErr != nil {
		c.JSON(403, "Error creating token")
		return
	}

	c.Writer.Header().Set("Authorization", token)

	c.JSON(200, foundUser.ToSimpleUser())
}

func(service *Service) CreateUser(c * gin.Context){
	var request types.CreateUserDto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First we have to check whether user with username already exists
	var foundUsers []db.User
	service.DB.Where("username = ?", request.Username).Find(&foundUsers).Limit(1)
	if len(foundUsers) > 0 {
		c.JSON(403, "This user already exists")
		return
	}

	hashedPassword, passwordErr := hashPassword(request.Password)

	if passwordErr != nil {
		c.JSON(400, gin.H{"message": "Error when hashing password"})
		c.Abort()
		return
	}

	user := db.User{
		Username: request.Username,
		Password: hashedPassword,
	}

	service.DB.Create(&user)

	token, tokenErr := CreateToken(user.ID)
	if tokenErr != nil {
		c.JSON(403, "Error creating token")
		return
	}

	c.Writer.Header().Set("Authorization", token)
	c.JSON(200, user.ToSimpleUser())
}
