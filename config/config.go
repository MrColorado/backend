package config

import "os"

type AwsConfigStruct struct {
	S3Location string
	S3UserName string
	S3Password string
}

type ScraperConfigStruct struct {
	S3Location string
	S3UserName string
	S3Password string
}

type Config struct {
	AwsConfig     AwsConfigStruct
	ScraperConfig ScraperConfigStruct
}

func GetConfig() Config {
	config := Config{
		AwsConfig: AwsConfigStruct{
			S3Location: os.Getenv("S3_LOCATION"),
			S3UserName: os.Getenv("S3_USERNAME"),
			S3Password: os.Getenv("S3_PASSWORD"),
		},
	}
	return config
}
