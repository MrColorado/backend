package awsWrapper

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/MrColorado/epubScraper/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	defaultRegion = "us-east-1"
	bucketName    = "novels"
)

type AwsClient struct {
	s3Client *s3.Client
}

func NewClient(awsConfig config.AwsConfigStruct) *AwsClient {
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               awsConfig.S3Location,
			SigningRegion:     defaultRegion,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:           defaultRegion,
		Credentials:      credentials.NewStaticCredentialsProvider(awsConfig.S3UserName, awsConfig.S3Password, ""),
		EndpointResolver: staticResolver,
	}

	return &AwsClient{
		s3Client: s3.NewFromConfig(cfg),
	}
}

func (client *AwsClient) UploadFile(filePath string, fileName string, content string) error {
	_, err := client.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", filePath, fileName)),
		Body:   strings.NewReader(content),
	})

	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			fileName, bucketName, filePath, err)
	}

	return err
}

func (client *AwsClient) DownLoadFile(filePath string) (string, error) {
	result, err := client.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, filePath, err)
		return "", err
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", filePath, err)
	}
	return string(body), err
}
