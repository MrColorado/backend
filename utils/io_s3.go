package utils

import (
	"encoding/json"
	"fmt"

	"github.com/MrColorado/epubScraper/awsWrapper"
)

type S3IO struct {
	awsClient *awsWrapper.AwsClient
}

func NewS3IO(awsClient *awsWrapper.AwsClient) S3IO {
	return S3IO{
		awsClient: awsClient,
	}
}

// ExportNovelChapter write novel chapter on s3
func (io S3IO) ExportNovelChapter(novelName string, novelChapterData NovelChapterData) error {
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
func (io S3IO) ExportMetaData(novelName string, novelMetaData NovelMetaData) error {
	content, err := json.Marshal(novelMetaData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshalize meta data of novel %s", novelMetaData.Title)
	}
	fmt.Printf("Export meta data of %s at path %s/meta_data.json\n", novelMetaData.Title, novelName)
	io.awsClient.UploadFile(novelName, "meta_data.json", content)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export meta data of novel %s", novelName)
	}
	return nil
}

// ImportNovelChapter read novel chapter from s3
func (io S3IO) ImportNovelChapter(novelName string, chapterData *NovelChapterData) error {
	content, err := io.awsClient.DownLoadFile(fmt.Sprintf("%s/raw", novelName), fmt.Sprintf("%04d.json", chapterData.Chapter))
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to get chapter %d of novel %s", chapterData.Chapter, novelName)
	}

	if json.Unmarshal([]byte(content), chapterData) != nil {
		return fmt.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}
	return nil
}

// ImportMetaData read novel meta data from s3
func (io S3IO) ImportMetaData(novelName string, novelMetaData *NovelMetaData) error {
	content, err := io.awsClient.DownLoadFile(novelName, "meta_data.json")
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to get meta_data of novel %s", novelName)
	}

	if json.Unmarshal([]byte(content), novelMetaData) != nil {
		return fmt.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}
	return nil
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
