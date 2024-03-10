package core

import (
	"encoding/json"
	"fmt"

	msgType "github.com/MrColorado/backend/internal/message"
	"github.com/MrColorado/backend/logger"
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
		return []byte{}, "", logger.Errorf("failed to find novel with id: %s", ID)
	}

	filepath := fmt.Sprintf("%s/epub", data.CoreData.Title)
	fileName := fmt.Sprintf("%s-%04d-%04d.epub", data.CoreData.Title, start, end)
	content, err := app.s3.DownLoadFile(filepath, fileName)
	if err != nil {
		return []byte{}, "", logger.Errorf("failed to get content of book %s/%s", filepath, fileName)
	}

	return content, data.CoreData.Title, nil
}

func (app *App) GetNovelByTitle(title string) (models.NovelData, error) {
	data, err := app.db.GetNovelByTitle(title)
	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel %s", title)
	}

	corverUrl, err := app.s3.GetPreSignedLink(data.CoreData.CoverPath)
	if err != nil {
		data.CoreData.CoverPath = ""
	} else {
		data.CoreData.CoverPath = corverUrl
	}

	return data, nil
}

func (app *App) GetNovelById(ID string) (models.NovelData, error) {
	data, err := app.db.GetNovelById(ID)
	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel with id %s", ID)
	}

	corverUrl, err := app.s3.GetPreSignedLink(data.CoreData.CoverPath)
	if err != nil {
		data.CoreData.CoverPath = ""
	} else {
		data.CoreData.CoverPath = corverUrl
	}

	return data, nil
}

func (app *App) ListNovels(startBy string) ([]models.PartialNovelData, error) {
	novels, err := app.db.ListNovels(startBy)
	if err != nil {
		return []models.PartialNovelData{}, logger.Errorf("failed to get list of novel")
	}

	for _, novel := range novels {
		corverUrl, err := app.s3.GetPreSignedLink(novel.CoverPath)
		if err != nil {
			continue
		}
		novel.CoverPath = corverUrl
	}

	return novels, nil
}

func (app *App) ListBook(ID string) ([]models.BookData, error) {
	books, err := app.db.ListBooks(ID)
	if err != nil {
		return []models.BookData{}, logger.Errorf("failed to find novel with id: %s", ID)
	}

	return books, nil
}

func (app *App) RequestNovel(title string) error {
	j, _ := json.Marshal(msgType.CanScrapeRqt{
		Title: title,
	})
	out, _ := json.Marshal(msgType.Message{
		Event:   "can_scrape",
		Payload: json.RawMessage(j),
	})

	resp, err := app.nats.Request("scrapable", out)
	if err != nil {
		return logger.Error("failed to request on nats")
	}
	var msg msgType.Message
	if json.Unmarshal(resp, &msg) != nil {
		return logger.Errorf("failed to unmarshal response : %s : %s", string(resp), err.Error())
	}

	var scrapeRsp msgType.CanScrapeRsp
	if json.Unmarshal(msg.Payload, &scrapeRsp) != nil {
		return logger.Errorf("failed to unmarshal response : %s : %s", string(resp), err.Error())
	}
	if len(scrapeRsp.ScraperName) == 0 {
		return logger.Errorf("can not scrape novel %s", title)
	}

	j, _ = json.Marshal(msgType.ScrapeNovelRqt{
		NovelTitle:  title,
		ScraperName: scrapeRsp.ScraperName,
	})
	out, _ = json.Marshal(msgType.Message{
		Event:   "scrape",
		Payload: json.RawMessage(j),
	})

	err = app.nats.PublishMsg(fmt.Sprintf("scraper.%s", scrapeRsp.ScraperName), out)
	if err != nil {
		return logger.Error("failed to publish msg")
	}
	return nil
}
