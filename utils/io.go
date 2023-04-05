package utils

// Scraper of each website should implement this interface
type IO interface {
	ExportNovelChapter(novelName string, novelChapterData NovelChapterData) error
	ExportMetaData(novelName string, novelMetaData NovelMetaData) error
	ImportMetaData(novelName string, novelMetaData *NovelMetaData) error
	NumberOfChapter(novelName string) (int, error)
	MataDataNotExist() bool
}
