package core

import (
	"encoding/json"
	"fmt"

	msgType "github.com/MrColorado/backend/internal/message"
	"github.com/MrColorado/backend/logger"
	"github.com/MrColorado/backend/server/internal/data"
	"github.com/MrColorado/backend/server/internal/models"
)

type App struct {
	s3   *data.S3Client
	db   *data.PostgresClient
	nats *data.NatsClient
}

func NewApp(s3 *data.S3Client, db *data.PostgresClient, nats *data.NatsClient) *App {
	return &App{
		s3:   s3,
		db:   db,
		nats: nats,
	}
}

func (app *App) GetBook(id string, start int, end int) ([]byte, string, error) {
	model, err := app.db.GetNovelByID(id)
	if err != nil {
		return []byte{}, "", logger.Errorf("failed to find novel with id: %s", id)
	}

	filepath := model.CoreData.Title + "/epub"
	fileName := fmt.Sprintf("%s-%04d-%04d.epub", model.CoreData.Title, start, end)
	content, err := app.s3.DownLoadFile(filepath, fileName)
	if err != nil {
		return []byte{}, "", logger.Errorf("failed to get content of book %s/%s", filepath, fileName)
	}

	return content, model.CoreData.Title, nil
}

func (app *App) GetNovelByTitle(title string) (models.NovelData, error) {
	model, err := app.db.GetNovelByTitle(title)
	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel %s", title)
	}

	corverURL, err := app.s3.GetPreSignedLink(model.CoreData.CoverPath)
	if err != nil {
		model.CoreData.CoverPath = ""
	} else {
		model.CoreData.CoverPath = corverURL
	}

	return model, nil
}

func (app *App) GetNovelByID(id string) (models.NovelData, error) {
	model, err := app.db.GetNovelByID(id)
	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel with id %s", id)
	}

	corverURL, err := app.s3.GetPreSignedLink(model.CoreData.CoverPath)
	if err != nil {
		model.CoreData.CoverPath = ""
	} else {
		model.CoreData.CoverPath = corverURL
	}

	return model, nil
}

func (app *App) ListNovels(startBy string) ([]models.PartialNovelData, error) {
	novels, err := app.db.ListNovels(startBy)
	if err != nil {
		return []models.PartialNovelData{}, logger.Errorf("failed to get list of novel")
	}

	for i := range novels {
		coverURL, err := app.s3.GetPreSignedLink(novels[i].CoverPath)
		if err != nil {
			continue
		}
		novels[i].CoverPath = coverURL
	}

	return novels, nil
}

func (app *App) ListBook(id string) ([]models.BookData, error) {
	books, err := app.db.ListBooks(id)
	if err != nil {
		return []models.BookData{}, logger.Errorf("failed to find novel with id: %s", id)
	}

	return books, nil
}

func (app *App) RequestNovel(title string) error {
	j, err := json.Marshal(msgType.CanScrapeRqt{
		Title: title,
	})
	if err != nil {
		return logger.Errorf("failed to marshal %+v", j)
	}
	out, err := json.Marshal(msgType.Message{
		Event:   "can_scrape",
		Payload: json.RawMessage(j),
	})
	if err != nil {
		return logger.Errorf("failed to marshal %+v", out)
	}

	resp, err := app.nats.Request("scrapable", out)
	if err != nil {
		return logger.Error("failed to request on nats")
	}
	var msg msgType.Message
	if err = json.Unmarshal(resp, &msg); err != nil {
		return logger.Errorf("failed to unmarshal response : %s : %s", string(resp), err.Error())
	}

	var scrapeRsp msgType.CanScrapeRsp
	if err = json.Unmarshal(msg.Payload, &scrapeRsp); err != nil {
		return logger.Errorf("failed to unmarshal response : %s : %s", string(resp), err.Error())
	}
	if len(scrapeRsp.ScraperName) == 0 {
		return logger.Errorf("can not scrape novel %s", title)
	}

	j, err = json.Marshal(msgType.ScrapeNovelRqt{
		NovelTitle:  title,
		ScraperName: scrapeRsp.ScraperName,
	})
	if err != nil {
		return logger.Errorf("failed to marshal %+v", j)
	}
	out, err = json.Marshal(msgType.Message{
		Event:   "scrape",
		Payload: json.RawMessage(j),
	})
	if err != nil {
		return logger.Errorf("failed to marshal %+v", out)
	}

	err = app.nats.PublishMsg("scraper."+scrapeRsp.ScraperName, out)
	if err != nil {
		return logger.Error("failed to publish msg")
	}
	return nil
}
