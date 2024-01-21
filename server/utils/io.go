package utils

import "github.com/MrColorado/backend/server/models"

// Scraper of each website should implement this interface
type IO interface {
	ExportNovelChapter(novelName string, novelChapterData models.NovelChapterData) error
	ExportMetaData(novelName string, novelMetaData models.NovelMetaData) error
	ImportNovelChapter(novelName string, chapter int) (models.NovelChapterData, error)
	ImportMetaData(novelName string) (models.NovelMetaData, error)
	ImportMetaDataById(novelId int) (models.NovelMetaData, error)
	NumberOfChapter(novelName string) (int, error)

	ExportBook(novelName string, bookName string, content []byte) error
	ListBook() ([]models.NovelMetaData, error)
}
