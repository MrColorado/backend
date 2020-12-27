package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func exportNovels(path string, novelName string, novelDatas []NovelData) {
	directoryPath := fmt.Sprintf("%s/%s", path, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if os.Mkdir(directoryPath, os.ModePerm) != nil {
			fmt.Printf("Failed to create directory : %s\n", directoryPath)
		}
	}
	for _, novelData := range novelDatas {
		j, err := json.Marshal(novelData)
		if err != nil {
			fmt.Printf("Failed to marshalize chapter %d of novel %s", novelData.Chapter, novelName)
			continue
		}
		fmt.Printf("%s/%04d.json\n", directoryPath, novelData.Chapter)
		ioutil.WriteFile(fmt.Sprintf("%s/%04d.json", directoryPath, novelData.Chapter), j, os.ModePerm)
	}
}

func exportNovel(path string, novelName string, novelData NovelData) {
	directoryPath := fmt.Sprintf("%s/%s", path, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if os.Mkdir(directoryPath, os.ModePerm) != nil {
			fmt.Printf("Failed to create directory : %s\n", directoryPath)
		}
	}
	j, err := json.Marshal(novelData)
	if err != nil {
		fmt.Printf("Failed to marshalize chapter %d of novel %s", novelData.Chapter, novelName)
		return
	}
	fmt.Printf("%s/%04d.json\n", directoryPath, novelData.Chapter)
	ioutil.WriteFile(fmt.Sprintf("%s/%04d.json", directoryPath, novelData.Chapter), j, os.ModePerm)
}
