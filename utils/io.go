package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ExportNovelChapter write novel chapter on disk
func ExportNovelChapter(path string, novelName string, novelChapterData NovelChapterData) {
	directoryPath := fmt.Sprintf("%s/%s", path, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if os.Mkdir(directoryPath, os.ModePerm) != nil {
			fmt.Printf("Failed to create directory : %s\n", directoryPath)
		}
	}
	j, err := json.Marshal(novelChapterData)
	if err != nil {
		fmt.Printf("Failed to marshalize chapter %d of novel %s\n", novelChapterData.Chapter, novelName)
		return
	}
	fmt.Printf("Export %s/%04d.json\n", directoryPath, novelChapterData.Chapter)
	ioutil.WriteFile(fmt.Sprintf("%s/%04d.json", directoryPath, novelChapterData.Chapter), j, os.ModePerm)
}

// ExportNovelMetaData write novel meta data on disk
func ExportNovelMetaData(path string, novelMetaData NovelMetaData) {
	directoryPath := fmt.Sprintf("%s/%s", path, novelMetaData.Title)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if os.Mkdir(directoryPath, os.ModePerm) != nil {
			fmt.Printf("Failed to create directory : %s\n", directoryPath)
		}
	}
	j, err := json.Marshal(novelMetaData)
	if err != nil {
		fmt.Printf("Failed to marshalize meta data of novel %s\n", novelMetaData.Title)
		return
	}
	fmt.Printf("Export meta data of %s at path %s/meta_data.json\n", novelMetaData.Title, directoryPath)
	ioutil.WriteFile(fmt.Sprintf("%s/meta_data.json", directoryPath), j, os.ModePerm)
}

// ImportNovel read novel chapter from disk
func ImportNovel(path string) (NovelChapterData, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return NovelChapterData{}, fmt.Errorf("Failed to read %s", path)
	}

	nodelData := NovelChapterData{}
	if json.Unmarshal(content, &nodelData) != nil {
		return NovelChapterData{}, fmt.Errorf("Failed to unmarshal %s", path)
	}

	return nodelData, nil
}

// ImportMetaData read novel meta data from disk
func ImportMetaData(path string) (NovelMetaData, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return NovelMetaData{}, fmt.Errorf("Failed to read %s", path)
	}

	novelMetaData := NovelMetaData{}
	if json.Unmarshal(content, &novelMetaData) != nil {
		return NovelMetaData{}, fmt.Errorf("Failed to unmarshal %s", path)
	}

	fmt.Printf("Import meta data from %s\n", path)
	return novelMetaData, nil
}

// NumberOfChapter return the chapter number of a novel
func NumberOfChapter(path string) int {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Printf("Failed to readDir %s\n", path)
		return 0
	}
	size := len(files)
	for _, file := range files {
		if file.Name() == "meta_data.json" || file.Name() == "cover" {
			size--
		}
	}
	return size
}
