package main

import (
	"github.com/MrColorado/backend/server/internal/config"
	"github.com/MrColorado/backend/server/internal/core"
	"github.com/MrColorado/backend/server/internal/dataHandler"
	"github.com/MrColorado/backend/server/internal/grpc"
)

func main() {
	cfg := config.GetConfig()
	app := core.NewApp(dataHandler.NewS3Client(cfg.AwsConfig), dataHandler.NewPostgresClient(cfg.PostgresConfig))
	server := grpc.NewSever(app)
	server.Run()
}
