//go:generate protoc --go_out=novelpb --go-grpc_out=novelpb novel.proto --experimental_allow_proto3_optional

package grpcWrapper

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/grpcWrapper/novelpb"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	io        utils.IO
	scraper   scraper.Scraper
	converter converter.Converter

	novelpb.UnimplementedNovelServerServer
}

func formatName(novelName string) string {
	return strings.TrimSpace(strings.ToLower(novelName))
}

func NewSever(io utils.IO, scraper scraper.Scraper, converter converter.Converter) *Server {
	return &Server{
		io:        io,
		scraper:   scraper,
		converter: converter,
	}
}

func (server *Server) GetNovel(ctx context.Context, req *novelpb.GetNovelRequest) (*novelpb.GetNovelResponse, error) {
	fmt.Println("Novel Service - Called GetNovel :", req.GetTitle())

	data, err := server.io.ImportMetaData(formatName(req.GetTitle()))
	if err != nil {
		return &novelpb.GetNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}
	return &novelpb.GetNovelResponse{
		Novel: &novelpb.NovelData{
			Title:       data.Title,
			Description: strings.Join(data.Summary, "\n"),
			Author:      data.Author,
			Chapter:     int64(data.NbChapter),
		},
	}, status.Error(codes.NotFound, "Novel was not found")
}

func (server *Server) ListNovel(ctx context.Context, req *novelpb.ListNovelRequest) (*novelpb.ListNovelResponse, error) {
	fmt.Println("Novel Service - Called ListNovel")

	datas, err := server.io.ListBooks()
	if err != nil {
		return &novelpb.ListNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	response := novelpb.ListNovelResponse{}
	for _, data := range datas {
		response.Novels = append(response.Novels, &novelpb.NovelData{
			Title:       data.Title,
			Description: strings.Join(data.Summary, "\n"),
			Author:      data.Author,
			Chapter:     int64(data.NbChapter),
		})
	}

	return &response, status.Error(codes.NotFound, "Not found")
}

func (server *Server) RequestNovel(_ *novelpb.NovelDemandRequest, _ novelpb.NovelServer_RequestNovelServer) error {
	return nil
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
