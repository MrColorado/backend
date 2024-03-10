//go:generate protoc --go_out=novelpb --go-grpc_out=novelpb ./protocol/novel.proto

package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/MrColorado/backend/logger"
	"github.com/MrColorado/backend/server/internal/core"
	"github.com/MrColorado/backend/server/internal/grpc/novelpb"
	"github.com/MrColorado/backend/server/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	app *core.App

	novelpb.UnimplementedNovelServerServer
}

func NewSever(app *core.App) *Server {
	return &Server{
		app: app,
	}
}

func (server *Server) GetBook(ctx context.Context, req *novelpb.GetBookRequest) (*novelpb.GetBookResponse, error) {
	logger.Infof("Novel Service - Called GetBook : ", req.GetNovelId())

	content, title, err := server.app.GetBook(req.NovelId, int(req.Chapter.Start), int(req.Chapter.End))
	if err != nil {
		return &novelpb.GetBookResponse{}, status.Error(codes.NotFound, "Not found")
	}

	response := novelpb.GetBookResponse{}
	response.Content = content
	response.Title = title
	return &response, nil
}

func (server *Server) GetNovel(ctx context.Context, req *novelpb.GetNovelRequest) (*novelpb.GetNovelResponse, error) {
	var err error
	var data models.NovelData

	switch req.OneofIdOrName.(type) {
	case *novelpb.GetNovelRequest_Id:
		logger.Infof("Novel Service - Called GetNovel : ", req.GetId())
		data, err = server.app.GetNovelById(req.GetId())
	case *novelpb.GetNovelRequest_Title:
		logger.Infof("Novel Service - Called GetNovel : ", req.GetTitle())
		data, err = server.app.GetNovelByTitle(req.GetTitle())
	}

	if err != nil {
		return &novelpb.GetNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	chpData, err := server.app.ListBook(data.CoreData.Id)
	if err != nil {
		return &novelpb.GetNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	var chapters []*novelpb.Chapter
	for _, chapter := range chpData {
		chapters = append(chapters, &novelpb.Chapter{
			Start: int64(chapter.Start),
			End:   int64(chapter.End),
		})
	}

	return &novelpb.GetNovelResponse{
		Novel: &novelpb.FullNovel{
			Novel: &novelpb.PartialNovel{
				Id:       data.CoreData.Id,
				Title:    data.CoreData.Title,
				Author:   data.CoreData.Author,
				Summary:  data.CoreData.Summary,
				CoverURL: data.CoreData.CoverPath,
				Genres:   data.CoreData.Genres,
			},
			NbChapter: int64(data.NbChapter),
			Tags:      data.Tags,
			Chapters:  chapters,
		},
	}, nil
}

func (server *Server) ListNovel(ctx context.Context, req *novelpb.ListNovelRequest) (*novelpb.ListNovelResponse, error) {
	logger.Info("Novel Service - Called ListNovel")

	datas, err := server.app.ListNovels(req.GetStartBy())
	if err != nil {
		return &novelpb.ListNovelResponse{}, status.Error(codes.NotFound, "Not found")
	}

	response := novelpb.ListNovelResponse{}
	for _, data := range datas {
		response.Novels = append(response.Novels, &novelpb.PartialNovel{
			Id:       data.Id,
			Title:    data.Title,
			Author:   data.Author,
			Summary:  data.Summary,
			CoverURL: data.CoverPath,
			Genres:   data.Genres,
		})
	}

	return &response, nil
}

func (server *Server) RequestNovel(ctx context.Context, req *novelpb.RequestNovelRequest) (*novelpb.RequestNovelResponse, error) {
	logger.Infof("Novel Service - Called RequestNovel : %s", req.GetTitle())
	err := server.app.RequestNovel(req.Title)
	if err != nil {
		return &novelpb.RequestNovelResponse{Success: false}, nil
	}
	return &novelpb.RequestNovelResponse{Success: true}, nil
}

func (server *Server) Run() {
	fmt.Println("Running novel Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55051")
	if err != nil {
		logger.Fatalf("Novel Service - ERROR:", err.Error())
	}

	s := grpc.NewServer()
	novelpb.RegisterNovelServerServer(s, server)

	logger.Infof("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		logger.Fatalf("Novel Service - ERROR:", err.Error())
	}
}
