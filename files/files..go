package files

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"recognizer/types"
)

func GetS3Client() *s3.Client {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Couldn't load default configuration. Here's why: %v\n", err)
		panic("Couldn't load default configuration.")
	}

	// Create S3 service client
	svc := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://fly.storage.tigris.dev")
		o.Region = "auto"
	})

	return svc
}

type Service struct {
	types.ServiceConfig
}

func NewFilesService(config types.ServiceConfig) Service {
	return Service{config}
}

func (service *Service) UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	key := uuid.New().String()

	_, err = service.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("recognizer"),
		Key:    aws.String(key),
		Body:   f,
	})

	c.JSON(200, gin.H{"url": key})
}
