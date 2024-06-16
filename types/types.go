package types

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

type ServiceConfig struct {
	DB *gorm.DB
	S3 *s3.Client
}
