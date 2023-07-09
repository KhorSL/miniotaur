package minio

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/khorsl/miniotaur/common/log/logger"
)

func NewMinioSvc() *MinioSvc {
	svc := &MinioSvc{}
	svc.initClient()
	return svc
}

func (m *MinioSvc) initClient() {
	logger := logger.NewLoggerWrapper(os.Getenv("DEFAULT_LOGGER_TYPE"), context.TODO())

	bucketUrl := os.Getenv("S3_BUCKET_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")
	bucketEndpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if bucketUrl != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           bucketUrl,
				SigningRegion: awsRegion,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(bucketEndpointResolver))

	if err != nil {
		logger.Error("Unable to load config", map[string]interface{}{
			"error": err,
		})
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	m.client = client
}

func (m *MinioSvc) ListBuckets(logger logger.Logger) ([]string, error) {
	buckets, err := m.client.ListBuckets(context.TODO(), nil)
	if err != nil {
		logger.Error("Unable to get buckets", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}

	var bucketNames []string
	for _, bucket := range buckets.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}

	return bucketNames, nil
}

func (m *MinioSvc) GetObject(bucket, key string, logger logger.Logger) ([]byte, error) {
	bucketName := bucket
	result, err := m.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Error("Couldn't get object", map[string]interface{}{
			"error":       err,
			"bucket_name": bucketName,
			"key":         key,
		})
		return nil, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		logger.Error("Couldn't read object body", map[string]interface{}{
			"error":       err,
			"bucket_name": bucketName,
			"key":         key,
		})
		return nil, err
	}

	return body, nil
}

func (m *MinioSvc) UploadObject(bucket, key string, file io.Reader, logger logger.Logger) (string, error) {
	uploader := manager.NewUploader(m.client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		logger.Error("Unable to upload", map[string]interface{}{
			"error":       err,
			"bucket_name": bucket,
			"key":         key,
		})
		return "", err
	}

	return *result.ETag, nil
}
