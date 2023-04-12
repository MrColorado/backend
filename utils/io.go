package utils

import "github.com/MrColorado/epubScraper/models"

// Scraper of each website should implement this interface
type IO interface {
	ExportNovelChapter(novelName string, novelChapterData models.NovelChapterData) error
	ExportMetaData(novelName string, novelMetaData models.NovelMetaData) error
	ImportNovelChapter(novelName string, chapter int) (models.NovelChapterData, error)
	ImportMetaData(novelName string) (models.NovelMetaData, error)
	NumberOfChapter(novelName string) (int, error)

	ExportBook(novelName string, bookName string, content []byte) error
	ListBooks() ([]models.NovelMetaData, error)
}
