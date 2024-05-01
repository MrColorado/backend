package main

import (
	"github.com/MrColorado/backend/logger"
	"github.com/MrColorado/backend/server/internal/config"
	"github.com/MrColorado/backend/server/internal/core"
	"github.com/MrColorado/backend/server/internal/data"
	"github.com/MrColorado/backend/server/internal/grpc"
)

func main() {
	config.InitLogger()
	cfg := config.GetConfig()
	nats, err := data.NewNatsClient(cfg.NatsConfig)
	if err != nil {
		logger.Info(err.Error())
		return
	}
	app := core.NewApp(data.NewS3Client(cfg.AwsConfig), data.NewPostgresClient(cfg.PostgresConfig), nats)
	server := grpc.NewSever(app)
	server.Run()
}

// func main() {
// 	config.InitLogger()
// 	config := config.GetConfig()

// 	db := data.NewPostgresClient(config.PostgresConfig)
// 	db.ListNovels()
// }
