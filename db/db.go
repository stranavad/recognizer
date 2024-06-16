package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Exam struct {
	BaseModel
	Name   string  `json:"name" binding:"required"`
	Groups []Group `gorm:"foreignKey:ExamId" json:"groups"`
	Items  []Item  `gorm:"foreignKey:ExamId" json:"items"`
}

type Group struct {
	BaseModel
	Name   string `json:"name" binding:"required"`
	ExamId uint   `json:"examId" binding:"required"`
	Items  []Item `gorm:"foreignKey:GroupId" json:"items"`
}

type Item struct {
	BaseModel
	Name    string `json:"name" binding:"required"`
	Image   string `json:"image"`
	GroupId uint   `json:"groupId"`
	ExamId  uint   `json:"examId"`
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

	err = db.AutoMigrate(&Exam{}, &Group{}, &Item{})
	if err != nil {
		panic("Failed to auto migrate")
	}

	return db
}
