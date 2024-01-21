package main

import (
	"github.com/MrColorado/backend/bookHandler/configuration"
	"github.com/MrColorado/backend/bookHandler/converter"
	"github.com/MrColorado/backend/bookHandler/dataWrapper"
	"github.com/MrColorado/backend/bookHandler/utils"
)

func main() {
	config := configuration.GetConfig()
	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

	io := utils.NewS3IO(awsClient, postgresClient)
	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

	conv.ConvertPartialNovel("big life", 1, 100)
}
