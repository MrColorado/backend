//go:generate protoc --go_out=novelpb --go-grpc_out=novelpb novel.proto

package grpcWrapper

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/grpcWrapper/novelpb"
	"github.com/MrColorado/epubScraper/models"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	io        utils.S3IO
	scraper   scraper.Scraper
	converter converter.Converter

	novelpb.UnimplementedNovelServerServer
}

func formatName(novelName string) string {
	return strings.TrimSpace(strings.ToLower(novelName))
}

func NewSever(io utils.S3IO, scraper scraper.Scraper, converter converter.Converter) *Server {
	return &Server{
		io:        io,
		scraper:   scraper,
		converter: converter,
	}
}

func (server *Server) GetNovel(ctx context.Context, req *novelpb.GetNovelRequest) (*novelpb.GetNovelResponse, error) {
	var err error
	var data models.NovelMetaData

	if req.GetId() != 0 {
		fmt.Println("Novel Service - Called GetNovel : ", req.GetId())
		data, err = server.io.ImportMetaDataById(int(req.GetId()))
	} else {
		fmt.Println("Novel Service - Called GetNovel : ", req.GetTitle())
		data, err = server.io.ImportMetaData(formatName(req.GetTitle()))
	}
	if err != nil {
		fmt.Println(err.Error())
		return &novelpb.GetNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	chaptersData, err := server.io.ListBooks(data.Id)
	if err != nil {
		fmt.Println(err.Error())
		return &novelpb.GetNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	var chapters []*novelpb.Chapter
	for _, chapter := range chaptersData {
		chapters = append(chapters, &novelpb.Chapter{
			Start: int64(chapter.Start),
			End:   int64(chapter.End),
		})
	}

	return &novelpb.GetNovelResponse{
		Novel: &novelpb.FullNovel{
			Novel: &novelpb.PartialNovel{
				Id:        int64(data.Id),
				Title:     data.Title,
				ImagePath: data.CoverPath,
			},
			Author:    data.Author,
			Summary:   strings.Join(data.Summary, "\n"),
			NbChapter: int64(data.NbChapter),
			Chapters:  chapters,
		},
	}, nil
}

func (server *Server) ListNovel(ctx context.Context, req *novelpb.ListNovelRequest) (*novelpb.ListNovelResponse, error) {
	fmt.Println("Novel Service - Called ListNovel")

	datas, err := server.io.ListNovels()
	if err != nil {
		fmt.Println(err.Error())
		return &novelpb.ListNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}
	fmt.Printf("Book size %d\n", len(datas))

	response := novelpb.ListNovelResponse{}
	for _, data := range datas {
		response.Novels = append(response.Novels, &novelpb.PartialNovel{
			Id:        int64(data.Id),
			Title:     data.Title,
			ImagePath: data.CoverPath,
		})
	}

	return &response, nil
}

func (server *Server) ScrapeNovel(ctx context.Context, req *novelpb.ScrapeNovelRequest) (*novelpb.ScrapeNovelResponse, error) {
	fmt.Println("Novel Service - Called RequestNovel")

	if !server.scraper.CanScrapeNovel(req.GetTitle()) {
		return &novelpb.ScrapeNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	go server.scraper.ScrapeNovel(req.GetTitle())
	return &novelpb.ScrapeNovelResponse{}, nil
}

func (server *Server) Run() {
	fmt.Println("Running novel Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55051")
	if err != nil {
		fmt.Println("Novel Service - ERROR:", err.Error())
	}

	s := grpc.NewServer()
	novelpb.RegisterNovelServerServer(s, server)

	fmt.Printf("Server started at %v\n", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		fmt.Println("Novel Service - ERROR:", err.Error())
	}

}
