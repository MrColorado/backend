//go:generate protoc --go_out=novelpb --go-grpc_out=novelpb novel.proto

package grpcWrapper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/file"
	"github.com/MrColorado/epubScraper/grpcWrapper/novelpb"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	chunkSize = 1048576 // 1 MB
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

func (server *Server) GetBook(req *novelpb.GetBookRequest, bookServer novelpb.NovelServer_GetBookServer) error {
	fmt.Println("Novel Service - Called GetBook : ", req.GetNovelId())

	data, err := server.io.ImportMetaDataById(req.NovelId)
	if err != nil {
		fmt.Println(err.Error())
		return status.Error(codes.NotFound, "Not found")
	}

	content, err := server.io.GetBook(data.Title, int(req.GetChapter().GetStart()), int(req.GetChapter().GetEnd()))
	if err != nil {
		fmt.Println(err.Error())
		return status.Error(codes.NotFound, "Not found")
	}

	f := file.NewFile(fmt.Sprintf("%s-%04d-%04d.epub", data.Title, req.GetChapter().GetStart(), req.GetChapter().GetEnd()), "epub", len(content), bytes.NewReader(content))
	err = bookServer.SendHeader(f.Metadata())
	if err != nil {
		return status.Error(codes.Internal, "error during sending header")
	}

	var n int
	chunk := &novelpb.GetBookResponse{Chunk: make([]byte, chunkSize)}

Loop:
	for {
		n, err = f.Read(chunk.Chunk)
		switch err {
		case nil:
		case io.EOF:
			break Loop
		default:
			return status.Errorf(codes.Internal, "io.ReadAll: %v", err)
		}
		chunk.Chunk = chunk.Chunk[:n]
		serverErr := bookServer.Send(chunk)
		if serverErr != nil {
			return status.Errorf(codes.Internal, "server.Send: %v", serverErr)
		}
	}

	return nil
}

func (server *Server) GetNovel(ctx context.Context, req *novelpb.GetNovelRequest) (*novelpb.GetNovelResponse, error) {
	fmt.Println("Novel Service - Called GetNovel : ", req.GetTitle())

	data, err := server.io.GetNovel(formatName(req.GetTitle()))
	if err != nil {
		fmt.Println(err.Error())
		return &novelpb.GetNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	chaptersData, err := server.io.ListBooks(data.CoreData.Title)
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
				Title:    data.CoreData.Title,
				Author:   data.CoreData.Author,
				Summary:  data.CoreData.Summary,
				CoverURL: "data.CoreData.CoverPath", // TODO GetPreSigned maybe not here
				Genres:   data.CoreData.Genres,
			},
			NbChapter: int64(data.NbChapter),
			Tags:      data.Tags,
			Chapters:  chapters,
		},
	}, nil
}

func (server *Server) ListNovel(ctx context.Context, req *novelpb.ListNovelRequest) (*novelpb.ListNovelResponse, error) {
	fmt.Println("Novel Service - Called ListNovel")

	datas, err := server.io.ListNovels(req.GetStartBy())
	if err != nil {
		fmt.Println(err.Error())
		return &novelpb.ListNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	response := novelpb.ListNovelResponse{}
	for _, data := range datas {
		response.Novels = append(response.Novels, &novelpb.PartialNovel{
			Title:    data.Title,
			Author:   data.Author,
			Summary:  data.Summary,
			CoverURL: "data.CoverPath", // TODO coverURL
			Genres:   data.Genres,
		})
	}

	return &response, nil
}

func (server *Server) RequestNovel(ctx context.Context, req *novelpb.RequestNovelRequest) (*novelpb.RequestNovelResponse, error) {
	fmt.Printf("Novel Service - Called RequestNovel : %s\n", req.GetTitle())

	if !server.scraper.CanScrapeNovel(req.GetTitle()) {
		fmt.Println("Faield to scrape novel")
		return &novelpb.RequestNovelResponse{Success: false}, status.Error(codes.NotFound, "Not found")
	}

	go server.scraper.ScrapeNovel(req.GetTitle())
	return &novelpb.RequestNovelResponse{Success: true}, nil
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
