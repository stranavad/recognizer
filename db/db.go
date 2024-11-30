package db

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Group struct {
	BaseModel
	Name   string `json:"name" binding:"required"`
	ExamID uint   `json:"examId" binding:"required"`
	Items  []Item `json:"items"`
	Exam Exam
}

type Exam struct {
	BaseModel
	Name   string  `json:"name" binding:"required"`
	Groups []Group `json:"groups"`
	Items  []Item  `json:"items"`
	UserID uint
}

type Item struct {
	BaseModel
	Name    string `json:"name" binding:"required"`
	Image   string `json:"image"`
	GroupID uint   `json:"groupId"`
	ExamID  uint   `json:"examId"`
	Exam Exam
}

type User struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex"`
	Password string
	Exams []Exam
}


type ScorePoint struct {
	BaseModel
	ExamID uint
	ItemID uint
	UserID uint
	Correct bool
}



func (user *User) ToSimpleUser() SimpleUser {
	return SimpleUser{
		ID:       user.ID,
		Username: user.Username,
	}
}

type SimpleUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func GetDB() *gorm.DB {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Println("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic("Failed to open database connection")
	}

	err = db.AutoMigrate(&Group{}, &Item{}, &User{}, &Exam{}, &ScorePoint{})
	if err != nil {
		panic("Failed to auto migrate")
	}

	return db
}
