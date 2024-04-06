package main

import (
	"github.com/MrColorado/backend/logger"
	"github.com/MrColorado/backend/server/internal/config"
	"github.com/MrColorado/backend/server/internal/core"
	"github.com/MrColorado/backend/server/internal/dataHandler"
	"github.com/MrColorado/backend/server/internal/grpc"
)

func main() {
	config.InitLogger()
	cfg := config.GetConfig()
	nats, err := dataHandler.NewNatsClient(cfg.NatsConfig)
	if err != nil {
		logger.Info(err.Error())
		return
	}
	app := core.NewApp(dataHandler.NewS3Client(cfg.AwsConfig), dataHandler.NewPostgresClient(cfg.PostgresConfig), nats)
	server := grpc.NewSever(app)
	server.Run()
}

// func main() {
// 	config.InitLogger()
// 	config := config.GetConfig()

// 	db := dataHandler.NewPostgresClient(config.PostgresConfig)
// 	db.ListNovels()
// }
