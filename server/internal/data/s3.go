package data

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
	bucketName = "novel"
)

type S3Client struct {
	s3Client        *s3.Client
	preSignedClient *s3.PresignClient
}

func getS3Client(host string, accessKey string, secretKey string) *s3.Client {
	customInternResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               host,
			HostnameImmutable: true,
		}, nil
	})

	conf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customInternResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return s3.NewFromConfig(conf)
}

func NewS3Client(conf cfg.AwsConfigStruct) *S3Client {
	return &S3Client{
		s3Client:        getS3Client(conf.S3InternHost, conf.S3UserName, conf.S3Password),
		preSignedClient: s3.NewPresignClient(getS3Client(conf.S3ExternHost, conf.S3UserName, conf.S3Password)),
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
	presignedURL, err := client.preSignedClient.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filePath),
		},
	)
	if err != nil {
		return "", logger.Errorf("Couldn't presigned file %v : %v\n", filePath, err)
	}
	return presignedURL.URL, nil
}
