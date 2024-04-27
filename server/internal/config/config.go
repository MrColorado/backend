package config

import (
	"os"

	"github.com/MrColorado/backend/logger"
)

type AwsConfigStruct struct {
	S3UserName   string
	S3Password   string
	S3InternHost string
	S3ExternHost string
}

type PostgresConfigStruct struct {
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
}

type NatsConfigStruct struct {
	NatsHOST string
}

type Config struct {
	AwsConfig      AwsConfigStruct
	PostgresConfig PostgresConfigStruct
	NatsConfig     NatsConfigStruct
}

func InitLogger() {
	logger.Init(
		logger.Configuration{
			AppName: "book-handler",
			Logger: &logger.LoggerConfiguration{
				DevLogs:    true,
				StackTrace: true,
			},
		},
	)
}

func GetConfig() Config {
	config := Config{
		AwsConfig: AwsConfigStruct{
			S3UserName:   os.Getenv("S3_USERNAME"),
			S3Password:   os.Getenv("S3_PASSWORD"),
			S3InternHost: os.Getenv("S3_INTERN_HOST"),
			S3ExternHost: os.Getenv("S3_EXTERN_HOST"),
		},
		PostgresConfig: PostgresConfigStruct{
			PostgresDB:       os.Getenv("POSTGRES_DB"),
			PostgresUser:     os.Getenv("POSTGRES_USER"),
			PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
			PostgresHost:     os.Getenv("POSTGRES_HOST"),
		}, NatsConfig: NatsConfigStruct{
			NatsHOST: os.Getenv("NATS_HOST"),
		},
	}
	return config
}
