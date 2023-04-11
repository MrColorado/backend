package grpcwrapper

import (
	"context"
	"log"
	"net"

	"github.com/MrColorado/epubScraper/novelpb"
	"google.golang.org/grpc"
)

type server struct {
	novelpb.UnimplementedNovelServerServer
}

func (*server) GetNovel(ctx context.Context, req *novelpb.GetNovelRequest) (*novelpb.GetNovelResponse, error) {
	log.Println("Novel Service - Called GetNovel :", req.NovelName)

	return &novelpb.GetNovelResponse{}, nil
}

func Test() {
	log.Println("Running Offer Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55051")
	if err != nil {
		log.Println("Offer Service - ERROR:", err.Error())
	}

	s := grpc.NewServer()
	novelpb.RegisterNovelServerServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("Offer Service - ERROR:", err.Error())
	}

}
