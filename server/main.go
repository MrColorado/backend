package main

import (
	"fmt"

	"github.com/MrColorado/backend/server/internal/config"
	"github.com/MrColorado/backend/server/internal/core"
	"github.com/MrColorado/backend/server/internal/dataHandler"
	"github.com/MrColorado/backend/server/internal/grpc"
)

func main() {
	fmt.Println("TEEEEST")
	cfg := config.GetConfig()
	nats, err := dataHandler.NewNatsClient(cfg.NatsConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	app := core.NewApp(dataHandler.NewS3Client(cfg.AwsConfig), dataHandler.NewPostgresClient(cfg.PostgresConfig), nats)
	server := grpc.NewSever(app)
	server.Run()
}
