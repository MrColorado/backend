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
		fmt.Printf("Failed to marshalize chapter %d of novel %s\n", novelChapterData.Chapter, novelName)
		return err
	}

	exportName := fmt.Sprintf("%04d.json", novelChapterData.Chapter)
	fmt.Printf("Export %s/%s\n", novelName, exportName)
	io.awsClient.UploadFile(novelName, exportName, content)
	return nil
}

// ExportMetaData write novel meta data on s3
func (io S3IO) ExportMetaData(novelName string, novelMetaData NovelMetaData) error {
	content, err := json.Marshal(novelMetaData)
	if err != nil {
		fmt.Printf("Failed to marshalize meta data of novel %s\n", novelMetaData.Title)
		return err
	}
	fmt.Printf("Export meta data of %s at path %s/meta_data.json\n", novelMetaData.Title, novelName)
	io.awsClient.UploadFile(novelName, "meta_data.json", content)
	return nil
}

// ImportMetaData read novel meta data from disk
func (io S3IO) ImportMetaData(novelName string, novelMetaData *NovelMetaData) error {
	content, err := io.awsClient.DownLoadFile(novelName, "meta_data.json")
	if err != nil {
		return fmt.Errorf("failed to get meta_data of novel %s", novelName)
	}

	if json.Unmarshal([]byte(content), novelMetaData) != nil {
		return fmt.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}

	return nil
}

// NumberOfChapter return the chapter number of a novel
func (io S3IO) NumberOfChapter(_ string) (int, error) {
	// files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", path, novelName))

	// if err != nil {
	// 	fmt.Printf("Failed to readDir %s\n", path)
	// 	return 0
	// }
	// size := len(files)
	// for _, file := range files {
	// 	if file.Name() == "meta_data.json" || file.Name() == "cover" {
	// 		size--
	// 	}
	// }
	// return size
	return 0, nil
}

// MataDataNotExist check if meta data are already exported
func (io S3IO) MataDataNotExist() bool {
	return true

	// directoryPath := fmt.Sprintf("%s", path)
	// if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
	// 	if os.Mkdir(directoryPath, os.ModePerm) != nil {
	// 		fmt.Printf("Failed to create directory : %s\n", directoryPath)
	// 	}
	// }

	// _, err := os.Stat(fmt.Sprintf("%s/meta_data.json", path))
	// return os.IsNotExist(err)
}
