package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DiskIO struct {
	outputPath string
}

func NewDiskIO(outputPath string) DiskIO {
	return DiskIO{
		outputPath: outputPath,
	}
}

// ExportNovelChapter write novel chapter on disk
func (io DiskIO) ExportNovelChapter(novelName string, chapterData NovelChapterData) error {
	directoryPath := fmt.Sprintf("%s/%s/raw", io.outputPath, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err = os.Mkdir(directoryPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf("failed to create directory : %s", directoryPath)
		}
	}
	j, err := json.Marshal(chapterData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshalize chapter %d of novel %s", chapterData.Chapter, novelName)
	}
	fmt.Printf("Export %s/%04d.json\n", directoryPath, chapterData.Chapter)
	err = ioutil.WriteFile(fmt.Sprintf("%s/%04d.json", directoryPath, chapterData.Chapter), j, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export chapter %d of novel %s", chapterData.Chapter, novelName)
	}
	return nil
}

// ExportMetaData write novel meta data on disk
func (io DiskIO) ExportMetaData(novelName string, novelMetaData NovelMetaData) error {
	directoryPath := fmt.Sprintf("%s/%s", io.outputPath, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf("failed to create directory : %s", directoryPath)
		}
	}
	j, err := json.Marshal(novelMetaData)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to marshalize meta data of novel %s", novelMetaData.Title)
	}
	fmt.Printf("Export meta data of novel %s at path %s/meta_data.json\n", novelMetaData.Title, directoryPath)
	err = ioutil.WriteFile(fmt.Sprintf("%s/meta_data.json", directoryPath), j, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export meta data of novel %s", novelName)
	}
	return nil
}

// ImportNovelChapter read novel chapter from disk
func (io DiskIO) ImportNovelChapter(novelName string, chapterData *NovelChapterData) error {
	directoryPath := fmt.Sprintf("%s/%s/raw", io.outputPath, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf("failed to create directory : %s", directoryPath)
		}
	}
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%04d.json", directoryPath, chapterData.Chapter))
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to get chapter %d of novel %s", chapterData.Chapter, novelName)
	}

	if json.Unmarshal(content, &chapterData) != nil {
		return fmt.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}

	return nil
}

// ImportMetaData read novel meta data from disk
func (io DiskIO) ImportMetaData(novelName string, novelMetaData *NovelMetaData) error {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/meta_data.json", io.outputPath, novelName))

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to get meta_data of novel %s", novelName)
	}

	if json.Unmarshal(content, &novelMetaData) != nil {
		return fmt.Errorf("failed to unmarshal metadata of novel %s", novelName)
	}

	return nil
}

// NumberOfChapter return the chapter number of a novel
func (io DiskIO) NumberOfChapter(novelName string) (int, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s/raw", io.outputPath, novelName))
	if err != nil {
		fmt.Println(err.Error())
		return 0, fmt.Errorf("failed to list files of novel %s", novelName)
	}
	return len(files), nil
}

// ExportBook return the chapter number of a novel
func (io DiskIO) ExportBook(novelName string, bookName string, content []byte) error {
	directoryPath := fmt.Sprintf("%s/%s/epub", io.outputPath, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf("failed to create directory : %s", directoryPath)
		}
	}
	exportName := fmt.Sprintf("%s.epub", bookName)
	fmt.Printf("Export book %s of novel %s\n", exportName, novelName)
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s", directoryPath, exportName), content, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export book %s of novel %s", exportName, novelName)
	}
	return nil
}
