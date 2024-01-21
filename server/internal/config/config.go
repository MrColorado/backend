package config

import "os"

type AwsConfigStruct struct {
	S3Location string
	S3UserName string
	S3Password string
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

func GetConfig() Config {
	config := Config{
		AwsConfig: AwsConfigStruct{
			S3Location: os.Getenv("S3_HOST"),
			S3UserName: os.Getenv("S3_USERNAME"),
			S3Password: os.Getenv("S3_PASSWORD"),
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
