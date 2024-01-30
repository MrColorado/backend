package dataStore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	cfg "github.com/MrColorado/backend/bookHandler/internal/config"
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

func NewAwsClient(awsConfig cfg.AwsConfigStruct) *S3Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
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
		fmt.Println(err)
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
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
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
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, filePath, err)
		return []byte{}, err
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", filePath, err)
	}
	return body, err
}
