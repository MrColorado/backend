package data

import (
	"bytes"
	"context"
	"fmt"
	"io"

	bkhConfig "github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucketName = "novel"
)

type S3Client struct {
	s3Client        *s3.Client
	preSignedClient *s3.PresignClient
}

func NewAwsClient(awsConfig bkhConfig.AwsConfigStruct) *S3Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               awsConfig.S3Location,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsConfig.S3UserName, awsConfig.S3Password, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		logger.Warnf("failed to create aws config : %s", err.Error())
	}

	s3Client := s3.NewFromConfig(cfg)
	return &S3Client{
		s3Client:        s3Client,
		preSignedClient: s3.NewPresignClient(s3Client),
	}
}

func (client *S3Client) UploadFile(filePath string, fileName string, content []byte) error {
	_, err := client.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", filePath, fileName)),
		Body:   bytes.NewReader(content),
	})

	if err != nil {
		return logger.Errorf("Couldn't upload file %v to %v:%v : %v\n",
			fileName, bucketName, filePath, err)
	}

	return err
}

func (client *S3Client) DownLoadFile(filePath string, fileName string) ([]byte, error) {
	result, err := client.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", filePath, fileName)),
	})
	if err != nil {
		return []byte{}, logger.Errorf("Couldn't get object %v:%v : %v\n", bucketName, filePath, err)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		logger.Errorf("Couldn't read object body from %v : %v\n", filePath, err)
	}
	return body, nil
}

func (client *S3Client) RemoveFile(filePath string, fileName string) error {
	_, err := client.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", filePath, fileName)),
	})
	if err != nil {
		return logger.Errorf("Failed to remove file %v:%v : %v\n", bucketName, filePath, err)
	}

	return nil
}
