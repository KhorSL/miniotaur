package minio

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MinioSvc struct {
	client *s3.Client
}
