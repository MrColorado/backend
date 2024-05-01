package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/MrColorado/backend/book-handler/internal/data"
	"github.com/MrColorado/backend/book-handler/internal/models"
	"github.com/MrColorado/backend/logger"
	"github.com/google/uuid"
)

const (
	coverDirectory = "covers"
)

type App struct {
	s3 *data.S3Client
	db *data.PostgresClient
}

func NewApp(s3 *data.S3Client, db *data.PostgresClient) *App {
	return &App{
		s3: s3,
		db: db,
	}
}

// ExportNovelChapter write novel chapter on s3
func (app *App) ExportNovelChapter(novelName string, novelChapterData models.NovelChapterData) error {
	content, err := json.Marshal(novelChapterData)
	if err != nil {
		return logger.Errorf("failed to marshalize chapter %d of novel %s", novelChapterData.Chapter, novelName)
	}

	exportName := fmt.Sprintf("%04d.json", novelChapterData.Chapter)
	logger.Infof("Export %s/raw/%s", novelName, exportName)
	err = app.s3.UploadFile(novelName+"/raw", exportName, content)
	if err != nil {
		return logger.Errorf("failed to export chapter %d of novel %s", novelChapterData.Chapter, novelName)
	}
	return nil
}

// ExportMetaData write novel meta data on s3
func (app *App) ExportMetaData(novelName string, model models.NovelMetaData, genre bool) error {
	coverName := "cover.jpg"
	model.CoverPath = fmt.Sprintf("%s/%s", novelName, coverName)
	err := app.s3.UploadFile(novelName, coverName, model.CoverData)
	if err != nil {
		return logger.Errorf("failed to save cover of novel %s in s3", model.Title)
	}
	model.CoverPath = fmt.Sprintf("%s/%s", novelName, coverName)

	if genre {
		for _, name := range model.Genres {
			err = app.db.InsertOrUpdateGenre(name)
			if err != nil {
				return logger.Errorf("failed to export genre %s of novel %s in database", name, model.Title)
			}
		}
	}

	err = app.db.InsertOrUpdateNovel(model, genre)
	if err != nil {
		return logger.Errorf("failed to export metedata of novel %s in database", model.Title)
	}

	return nil
}

// ExportBook return the chapter number of a novel
func (app *App) ExportBook(novelName string, bookName string, content []byte, metaData models.BookData) error {
	filePath := novelName + "/epub"

	book, err := app.db.GetBookByTitle(novelName, metaData.Start)
	if err == nil && metaData.End > book.End {
		logger.Info("New book contains more chapters new %d, current %s removing older book from S3", metaData.End, book.End)
		fileName := fmt.Sprintf("%s-%04d-%04d.epub", novelName, metaData.Start, metaData.End)
		if err = app.s3.RemoveFile(filePath, fileName); err != nil {
			logger.Errorf("Failed to remove file %s at %s", fileName, filePath)
		}
	}

	novel, err := app.db.GetNovelByTitle(novelName)
	if err != nil {
		return logger.Errorf("Failed to get novel %s : %s", novelName, err.Error())
	}

	fileName := bookName + ".epub"
	err = app.s3.UploadFile(filePath, fileName, content)
	if err != nil {
		return logger.Errorf("failed to export book %s of novel %s", fileName, novelName)
	}
	logger.Infof("Export book %s of novel %s", fileName, novelName)

	metaData.NovelID = novel.CoreData.ID
	err = app.db.InsertOrUpdateBook(metaData)
	if err != nil {
		return logger.Errorf("failed to save book %s", novelName)
	}

	return nil
}

// GetNovelChapter read novel chapter from s3
func (app *App) GetNovelChapter(novelName string, chapter int) (models.NovelChapterData, error) {
	content, err := app.s3.DownLoadFile(novelName+"/raw", fmt.Sprintf("%04d.json", chapter))
	if err != nil {
		return models.NovelChapterData{}, logger.Errorf("failed to get chapter %d of novel %s", chapter, novelName)
	}

	model := models.NovelChapterData{}
	if json.Unmarshal(content, &model) != nil {
		return models.NovelChapterData{}, logger.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}
	return model, nil
}

func (app *App) GetMetaData(title string) (models.NovelMetaData, error) {
	model, err := app.db.GetNovelByTitle(title)
	if err != nil {
		return models.NovelMetaData{}, logger.Errorf("failed to get meta_data of novel %s", title)
	}
	return models.NovelToMeta(model), nil
}

func (app *App) GetNbChapter(title string) (int, error) {
	model, err := app.db.GetNovelByTitle(title)
	if err != nil {
		return 0, logger.Errorf("failed to get nb chapter of %s", title)
	}
	return model.NbChapter, nil
}

func (app *App) GetNovelByTitle(novelName string) (models.NovelData, error) {
	model, err := app.db.GetNovelByTitle(novelName)
	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get meta_data of novel %s", novelName)
	}
	return model, nil
}

func (app *App) GetCoverDiskPath(title string) (string, error) {
	fileName := "cover.jpg"
	buf, err := app.s3.DownLoadFile(title, fileName)
	if err != nil {
		return "", logger.Errorf("failed to get cover of novel %s", title)
	}
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return "", logger.Errorf("failed to decode data of novel %s's cover", title)
	}

	ok, err := exists(coverDirectory)
	if !ok || err != nil {
		os.Mkdir(coverDirectory, os.ModePerm)
	}

	// TODO check how bypass writing on disk
	path := fmt.Sprintf("%s/%s", coverDirectory, uuid.New().String())
	fo, err := os.Create(path)
	if err != nil {
		return "", logger.Errorf("failed to save cover of novel %s on path %s", title, path)
	}
	defer fo.Close()

	if err = jpeg.Encode(fo, img, nil); err != nil {
		return "", logger.Errorf("failed to save cover of novel %s on path %s", title, path)
	}
	for {
		n, err := fo.Read(buf)
		if err != nil && err != io.EOF {
			return "", logger.Errorf("failed to save cover of novel %s on path %s", title, path)
		}
		if n == 0 {
			break
		}

		if _, err := fo.Write(buf[:n]); err != nil {
			return "", logger.Errorf("failed to save cover of novel %s on path %s", title, path)
		}
	}

	return path, nil
}

func (app *App) RemoveCoverDiskPath(filepath string) {
	ok, err := exists(coverDirectory)
	if !ok || err != nil {
		return
	}
	err = os.Remove(filepath)
	if err != nil {
		logger.Warn("failed to remove %s : %s", filepath, err.Error())
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
