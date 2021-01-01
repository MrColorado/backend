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
		fmt.Printf("Failed to marshalize chapter %d of novel %s", novelChapterData.Chapter, novelName)
		return
	}
	fmt.Printf("Export %s/%04d.json\n", directoryPath, novelChapterData.Chapter)
	ioutil.WriteFile(fmt.Sprintf("%s/%04d.json", directoryPath, novelChapterData.Chapter), j, os.ModePerm)
}

// InportNovel read novel chapter on disk
func InportNovel(path string) (NovelChapterData, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return NovelChapterData{}, fmt.Errorf("Failed to read %s", path)
	}

	nodelData := NovelChapterData{}
	if json.Unmarshal(content, &nodelData) != nil {
		return NovelChapterData{}, fmt.Errorf("Failed to unmarshal %s", path)
	}

	fmt.Printf("Import %s\n", path)
	return nodelData, nil
}

// NumberOfChapter return the chapter number of a novel
func NumberOfChapter(path string) int {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Printf("Failed to readDir %s\n", path)
		return 0
	}
	return len(files)
}
