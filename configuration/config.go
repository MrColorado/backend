package configuration

import "os"

type AwsConfigStruct struct {
	S3Location string
	S3UserName string
	S3Password string
}

type PostgresConfigStruct struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

type ScraperConfigStruct struct {
}

type ConverterConfigStruct struct {
}

type Config struct {
	AwsConfig       AwsConfigStruct
	ScraperConfig   ScraperConfigStruct
	PostgresConfig  PostgresConfigStruct
	ConverterConfig ConverterConfigStruct
}

func GetConfig() Config {
	config := Config{
		AwsConfig: AwsConfigStruct{
			S3Location: os.Getenv("S3_LOCATION"),
			S3UserName: os.Getenv("S3_USERNAME"),
			S3Password: os.Getenv("S3_PASSWORD"),
		},
		PostgresConfig: PostgresConfigStruct{
			PostgresUser:     os.Getenv("POSTGRES_USER"),
			PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
			PostgresDB:       os.Getenv("POSTGRES_DB"),
		},
	}
	return config
}
