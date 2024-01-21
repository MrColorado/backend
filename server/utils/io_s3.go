package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/MrColorado/backend/server/dataWrapper"
	"github.com/MrColorado/backend/server/models"
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
func (s3io S3IO) ExportNovelChapter(novelName string, novelChapterData models.NovelChapterData) error {
	content, err := json.Marshal(novelChapterData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshalize chapter %d of novel %s", novelChapterData.Chapter, novelName)
	}

	exportName := fmt.Sprintf("%04d.json", novelChapterData.Chapter)
	fmt.Printf("Export %s/raw/%s\n", novelName, exportName)
	s3io.awsClient.UploadFile(fmt.Sprintf("%s/raw", novelName), exportName, content)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export chapter %d of novel %s", novelChapterData.Chapter, novelName)
	}
	return nil
}

// ExportMetaData write novel meta data on s3
func (s3io S3IO) ExportMetaData(novelName string, data models.NovelMetaData) error {
	coverName := "cover.jpg"
	data.CoverPath = fmt.Sprintf("%s/%s", novelName, coverName)
	err := s3io.awsClient.UploadFile(novelName, coverName, data.CoverData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to save cover of novel %s in s3", data.Title)
	}
	data.CoverPath = fmt.Sprintf("%s/%s", novelName, coverName)
	err = s3io.dbClient.InsertOrUpdateNovel(data)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export metedata of novel %s in database", data.Title)
	}
	return nil
}

// ExportBook return the chapter number of a novel
func (s3io S3IO) ExportBook(novelName string, bookName string, content []byte, metaData models.BookData) error {
	exportName := fmt.Sprintf("%s.epub", bookName)
	fmt.Printf("Export book %s of novel %s\n", exportName, novelName)
	err := s3io.awsClient.UploadFile(fmt.Sprintf("%s/epub", novelName), exportName, content)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export book %s of novel %s", exportName, novelName)
	}

	// TODO create data struct that contain every field instead of doing this kind on request
	novelData, err := s3io.dbClient.GetNovelByTitle(novelName)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to get novel %s", novelName)
	}

	metaData.NovelId = novelData.CoreData.Id
	err = s3io.dbClient.InsertOrUpdateBook(metaData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to save book %s", novelName)
	}

	return nil
}

// ImportNovelChapter read novel chapter from s3
func (s3io S3IO) ImportNovelChapter(novelName string, chapter int) (models.NovelChapterData, error) {
	content, err := s3io.awsClient.DownLoadFile(fmt.Sprintf("%s/raw", novelName), fmt.Sprintf("%04d.json", chapter))
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
func (s3io S3IO) ImportMetaData(novelName string) (models.NovelMetaData, error) {
	data, err := s3io.dbClient.GetNovelByTitle(novelName)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelMetaData{}, fmt.Errorf("failed to get meta_data of novel %s", novelName)
	}
	return models.NovelToMeta(data), nil
}

// ImportMetaData read novel meta data from s3
func (s3io S3IO) ImportMetaDataById(novelId string) (models.NovelMetaData, error) {
	data, err := s3io.dbClient.GetNovelById(novelId)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelMetaData{}, fmt.Errorf("failed to get meta_data of novel %s", novelId)
	}
	return models.NovelToMeta(data), nil
}

// NumberOfChapter return the chapter number of a novel
func (s3io S3IO) NumberOfChapter(novelName string) (int, error) {
	filesName, err := s3io.awsClient.ListFiles(fmt.Sprintf("%s/raw", novelName))
	if err != nil {
		fmt.Println(err.Error())
		return 0, fmt.Errorf("failed to list files of novel %s", novelName)
	}
	return len(filesName), nil
}

func (s3io S3IO) GetNovelById(id string) (models.NovelData, error) {
	data, err := s3io.dbClient.GetNovelById(id)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelData{}, fmt.Errorf("failed to get meta_data of novel %s", id)
	}
	return data, nil
}

func (s3io S3IO) GetNovelByTitle(novelName string) (models.NovelData, error) {
	data, err := s3io.dbClient.GetNovelByTitle(novelName)
	if err != nil {
		fmt.Println(err.Error())
		return models.NovelData{}, fmt.Errorf("failed to get meta_data of novel %s", novelName)
	}
	return data, nil
}

func (s3io S3IO) ListNovels(startBy string) ([]models.PartialNovelData, error) {
	novels, err := s3io.dbClient.ListNovels(startBy)
	if err != nil {
		fmt.Println(err)
		return []models.PartialNovelData{}, fmt.Errorf("failed to get list of novel")
	}

	for _, novel := range novels {
		corverUrl, err := s3io.awsClient.GetPreSignedLink(novel.CoverPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(corverUrl)
		novel.CoverPath = corverUrl
	}

	return novels, nil
}

func (s3io S3IO) ListBooks(novelName string) ([]models.BookData, error) {
	datas, err := s3io.dbClient.ListBooks(novelName)
	if err != nil {
		fmt.Println(err)
		return []models.BookData{}, fmt.Errorf("failed to get list of books")
	}

	return datas, nil
}

func (s3io S3IO) GetBook(novelName string, start int, end int) ([]byte, error) {
	filepath := fmt.Sprintf("%s/epub", novelName)
	bookName := fmt.Sprintf("%s-%04d-%04d.epub", novelName, start, end)
	content, err := s3io.awsClient.DownLoadFile(filepath, bookName)
	if err != nil {
		fmt.Println(err)
		return []byte{}, fmt.Errorf("failed to dowload book : %s", bookName)
	}
	return content, nil
}

func (s3io S3IO) GetCoverUrl(title string) (string, error) {
	novel, err := s3io.dbClient.GetNovelByTitle(title)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("failed to get cover url of novel %s novel", title)
	}
	corverUrl, err := s3io.awsClient.GetPreSignedLink(novel.CoreData.CoverPath)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("failed to get cover url of novel %s novel", title)
	}
	return corverUrl, nil
}

func (s3io S3IO) GetCover(title string) (string, error) {
	buf, err := s3io.awsClient.DownLoadFile(title, "cover.jpg")
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("1 failed to get cover of novel %s", title)
	}
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("1 failed to get cover of novel %s", title)
	}

	path := "./cover.jpg"
	fo, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("2 failed to save cover of novel %s", title)
	}
	defer fo.Close()
	if err = jpeg.Encode(fo, img, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	for {
		// read a chunk
		n, err := fo.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("failed to save cover of novel %s : error %s", title, err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			fmt.Printf("failed to save cover of novel %s : error %s", title, err)
		}
	}

	return path, nil
}
