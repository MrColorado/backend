package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MrColorado/epubScraper/dataWrapper"
	"github.com/MrColorado/epubScraper/models"
)

type S3IO struct {
	awsClient *dataWrapper.AwsClient
	dbClient  *dataWrapper.PostgresClient
}

func NewS3IO(awsClient *dataWrapper.AwsClient, dbClient *dataWrapper.PostgresClient) S3IO {
	return S3IO{
		awsClient: awsClient,
		dbClient:  dbClient,
	}
}

// ExportNovelChapter write novel chapter on s3
func (io S3IO) ExportNovelChapter(novelName string, novelChapterData models.NovelChapterData) error {
	content, err := json.Marshal(novelChapterData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshalize chapter %d of novel %s", novelChapterData.Chapter, novelName)
	}

	exportName := fmt.Sprintf("%04d.json", novelChapterData.Chapter)
	fmt.Printf("Export %s/raw/%s\n", novelName, exportName)
	io.awsClient.UploadFile(fmt.Sprintf("%s/raw", novelName), exportName, content)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export chapter %d of novel %s", novelChapterData.Chapter, novelName)
	}
	return nil
}

// ExportMetaData write novel meta data on s3
func (io S3IO) ExportMetaData(novelName string, data *models.NovelMetaData) error {
	data.Title = strings.ToLower(data.Title)

	err := io.dbClient.InsertOrUpdate(data)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export metedata of novel %s in database", data.Title)
	}
	return nil
}

// ImportNovelChapter read novel chapter from s3
func (io S3IO) ImportNovelChapter(novelName string, chapter int) (models.NovelChapterData, error) {
	content, err := io.awsClient.DownLoadFile(fmt.Sprintf("%s/raw", novelName), fmt.Sprintf("%04d.json", chapter))
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelChapterData{}, fmt.Errorf("failed to get chapter %d of novel %s", chapter, novelName)
	}

	data := models.NovelChapterData{}
	if json.Unmarshal([]byte(content), &data) != nil {
		return models.NovelChapterData{}, fmt.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}
	return data, nil
}

// ImportMetaData read novel meta data from s3
func (io S3IO) ImportMetaData(novelName string) (models.NovelMetaData, error) {
	data, err := io.dbClient.GetTitle(novelName)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelMetaData{}, fmt.Errorf("failed to get meta_data of novel %s", novelName)
	}
	return data, nil
}

// NumberOfChapter return the chapter number of a novel
func (io S3IO) NumberOfChapter(novelName string) (int, error) {
	filesName, err := io.awsClient.ListFiles(fmt.Sprintf("%s/raw", novelName))
	if err != nil {
		fmt.Println(err.Error())
		return 0, fmt.Errorf("failed to list files of novel %s", novelName)
	}
	return len(filesName), nil
}

// ExportBook return the chapter number of a novel
func (io S3IO) ExportBook(novelName string, bookName string, content []byte) error {
	exportName := fmt.Sprintf("%s.epub", bookName)
	fmt.Printf("Export book %s of novel %s\n", exportName, novelName)
	err := io.awsClient.UploadFile(fmt.Sprintf("%s/epub", novelName), exportName, content)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export book %s of novel %s", exportName, novelName)
	}
	return nil
}

func (io S3IO) ListBooks() ([]models.NovelMetaData, error) {
	datas, err := io.dbClient.List()
	if err != nil {
		fmt.Println(err)
		return []models.NovelMetaData{}, fmt.Errorf("failed to get list of novel")
	}

	return datas, nil
}
