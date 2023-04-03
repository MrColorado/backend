package utils

// Scraper of each website should implement this interface
type IO interface {
	ExportNovelChapter(path string, novelName string, novelChapterData NovelChapterData)
	ExportMetaData(novelName string, novelMetaData NovelMetaData)
	ImportMetaData(path string, novelName string) (NovelMetaData, error)
	NumberOfChapter(path string, novelName string) int
	MataDataNotExist(path string) bool
}
