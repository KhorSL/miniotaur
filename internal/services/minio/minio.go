package minio

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewMinioSvc() *MinioSvc {
	svc := &MinioSvc{}
	svc.initClient()
	return svc
}

func (m *MinioSvc) initClient() {
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
		log.Printf("Unable to load config: %v\n", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	m.Client = client
}

func (m *MinioSvc) ListBuckets() ([]string, error) {
	buckets, err := m.Client.ListBuckets(context.TODO(), nil)
	if err != nil {
		log.Printf("Unable to get buckets: %v\n", err)
		return nil, err
	}

	var bucketNames []string
	for _, bucket := range buckets.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}

	return bucketNames, nil
}

func (m *MinioSvc) GetObject(bucket, key string) ([]byte, error) {
	bucketName := bucket
	result, err := m.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, key, err)
		return nil, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", key, err)
		return nil, err
	}

	return body, nil
}

func (m *MinioSvc) UploadObject(bucket, key string, file io.Reader) (string, error) {
	uploader := manager.NewUploader(m.Client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Printf("Unable to upload: %v\n", err)
		return "", err
	}

	return *result.ETag, nil
}
