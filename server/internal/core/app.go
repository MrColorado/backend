package core

import (
	"encoding/json"
	"fmt"

	msgType "github.com/MrColorado/backend/internal/message"
	"github.com/MrColorado/backend/server/internal/dataHandler"
	"github.com/MrColorado/backend/server/internal/models"
)

type App struct {
	s3   *dataHandler.S3Client
	db   *dataHandler.PostgresClient
	nats *dataHandler.NatsClient
}

func NewApp(s3 *dataHandler.S3Client, db *dataHandler.PostgresClient, nats *dataHandler.NatsClient) *App {
	return &App{
		s3:   s3,
		db:   db,
		nats: nats,
	}
}

func (app *App) GetBook(ID string, start int, end int) ([]byte, string, error) {
	data, err := app.db.GetNovelById(ID)
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, "", fmt.Errorf("failed to find novel with id: %s", ID)
	}

	filepath := fmt.Sprintf("%s/epub", data.CoreData.Title)
	fileName := fmt.Sprintf("%s-%04d-%04d.epub", data.CoreData.Title, start, end)
	content, err := app.s3.DownLoadFile(filepath, fileName)
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, "", fmt.Errorf("failed to get content of book %s/%s", filepath, fileName)
	}

	return content, data.CoreData.Title, nil
}

func (app *App) GetNovelByTitle(title string) (models.NovelData, error) {
	data, err := app.db.GetNovelByTitle(title)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelData{}, fmt.Errorf("failed to get novel %s", title)
	}

	corverUrl, err := app.s3.GetPreSignedLink(data.CoreData.CoverPath)
	if err != nil {
		fmt.Println(err)
		data.CoreData.CoverPath = ""
	} else {
		data.CoreData.CoverPath = corverUrl
	}

	return data, nil
}

func (app *App) GetNovelById(ID string) (models.NovelData, error) {
	data, err := app.db.GetNovelById(ID)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelData{}, fmt.Errorf("failed to get novel with id %s", ID)
	}

	corverUrl, err := app.s3.GetPreSignedLink(data.CoreData.CoverPath)
	if err != nil {
		fmt.Println(err)
		data.CoreData.CoverPath = ""
	} else {
		data.CoreData.CoverPath = corverUrl
	}

	return data, nil
}

func (app *App) ListNovels(startBy string) ([]models.PartialNovelData, error) {
	novels, err := app.db.ListNovels(startBy)
	if err != nil {
		fmt.Println(err)
		return []models.PartialNovelData{}, fmt.Errorf("failed to get list of novel")
	}

	for _, novel := range novels {
		corverUrl, err := app.s3.GetPreSignedLink(novel.CoverPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		novel.CoverPath = corverUrl
	}

	return novels, nil
}

func (app *App) ListBook(ID string) ([]models.BookData, error) {
	books, err := app.db.ListBooks(ID)
	if err != nil {
		fmt.Println(err.Error())
		return []models.BookData{}, fmt.Errorf("failed to find novel with id: %s", ID)
	}

	return books, nil
}

func (app *App) RequestNovel(title string) error {
	fmt.Println("A")
	out, err := json.Marshal(msgType.Message{
		Event: "can_scrape",
		Payload: msgType.CanScrapeRqt{
			Title: title,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshal request")
	}

	fmt.Println("B")
	resp, err := app.nats.Request("scrapable", out)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to request on nats")
	}

	fmt.Println("C")
	var in msgType.CanScrapeRsp
	err = json.Unmarshal(resp, &in)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to unmarshal response")
	}

	if len(in.ScraperName) == 0 {
		return fmt.Errorf("can not scrape novel %s", title)
	}

	out, err = json.Marshal(msgType.Message{
		Event: "scrape",
		Payload: msgType.ScrapeNovelRqt{
			NovelTitle: title,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshal request")
	}

	fmt.Println("D")
	err = app.nats.PublishMsg(fmt.Sprintf("scraper:%s", in.ScraperName), out)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to publish msg")
	}
	fmt.Println("E")
	return nil
}
