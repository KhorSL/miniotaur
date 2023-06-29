package minio

import "github.com/aws/aws-sdk-go-v2/service/s3"

type MinioSvc struct {
	Client *s3.Client
}
