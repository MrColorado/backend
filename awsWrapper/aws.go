package awsWrapper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/MrColorado/epubScraper/configuration"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

func NewClient(awsConfig configuration.AwsConfigStruct) *AwsClient {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               awsConfig.S3Location,
			SigningRegion:     defaultRegion,
			HostnameImmutable: true,
		}, nil
	})

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsConfig.S3UserName, awsConfig.S3Password, "")),
	)

	return &AwsClient{
		s3Client: s3.NewFromConfig(cfg),
	}
}

func (client *AwsClient) UploadFile(filePath string, fileName string, content []byte) error {
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

func (client *AwsClient) DownLoadFile(filePath string, fileName string) ([]byte, error) {
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

func (client *AwsClient) ListFiles(filePath string) ([]string, error) {
	result, err := client.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
		return []string{}, nil
	}
	var filesName []string
	for _, data := range result.Contents {
		filesName = append(filesName, *data.Key)
	}
	return filesName, nil
}
