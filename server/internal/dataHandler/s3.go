package dataHandler

import (
	"context"
	"fmt"
	"io"

	"github.com/MrColorado/backend/logger"
	cfg "github.com/MrColorado/backend/server/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucketName = "novels"
)

type S3Client struct {
	s3Client        *s3.Client
	preSignedClient *s3.PresignClient
}

func NewS3Client(c cfg.AwsConfigStruct) *S3Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               c.S3Location,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.S3UserName, c.S3Password, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	s3Client := s3.NewFromConfig(cfg)
	return &S3Client{
		s3Client:        s3Client,
		preSignedClient: s3.NewPresignClient(s3Client),
	}
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
	return body, err
}

func (client *S3Client) GetPreSignedLink(filePath string) (string, error) {
	presignedUrl, err := client.preSignedClient.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filePath),
		},
	)
	if err != nil {
		return "", logger.Errorf("Couldn't presigned file %v : %v\n", filePath, err)
	}
	return presignedUrl.URL, nil
}
