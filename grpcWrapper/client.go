package grpcWrapper

import (
	"context"
	"fmt"
	"net"

	"github.com/MrColorado/epubScraper/novelpb"
	"github.com/MrColorado/epubScraper/utils"
	"google.golang.org/grpc"
)

type Server struct {
	io utils.IO

	novelpb.UnimplementedNovelServerServer
}

func (server *Server) GetNovel(ctx context.Context, req *novelpb.GetNovelRequest) (*novelpb.GetNovelResponse, error) {
	fmt.Println("Novel Service - Called GetNovel :", req.NovelName)

	return &novelpb.GetNovelResponse{}, nil
}

func (server *Server) ListNovel(ctx context.Context, req *novelpb.ListNovelRequest) (*novelpb.ListNovelResponse, error) {
	fmt.Println("Novel Service - Called ListNovel")

	return &novelpb.ListNovelResponse{}, nil
}

func NewSever(io utils.IO) *Server {
	return &Server{
		io: io,
	}
}

func (server *Server) Run() {
	fmt.Println("Running novel Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55051")
	if err != nil {
		fmt.Println("Novel Service - ERROR:", err.Error())
	}

	s := grpc.NewServer()
	novelpb.RegisterNovelServerServer(s, server)

	fmt.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		fmt.Println("Novel Service - ERROR:", err.Error())
	}

}
