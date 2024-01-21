package main

import (
	"github.com/MrColorado/backend/server/internal/config"
	"github.com/MrColorado/backend/server/internal/core"
	"github.com/MrColorado/backend/server/internal/dataStore"
	"github.com/MrColorado/backend/server/internal/grpc"
)

func main() {
	cfg := config.GetConfig()
	app := core.NewApp(dataStore.NewS3Client(cfg.AwsConfig), dataStore.NewPostgresClient(cfg.PostgresConfig))
	server := grpc.NewSever(app)
	server.Run()
}
