package main

import (
	"github.com/MrColorado/backend/server/configuration"
	"github.com/MrColorado/backend/server/dataWrapper"
	"github.com/MrColorado/backend/server/grpcWrapper"
	"github.com/MrColorado/backend/server/utils"
)

func main() {
	config := configuration.GetConfig()
	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

	io := utils.NewS3IO(awsClient, postgresClient)

	server := grpcWrapper.NewSever(io)
	server.Run()
}
