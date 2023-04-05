package utils

// Scraper of each website should implement this interface
type IO interface {
	ExportNovelChapter(novelName string, novelChapterData NovelChapterData) error
	ExportMetaData(novelName string, novelMetaData NovelMetaData) error
	ImportNovelChapter(novelName string, novelChapterData *NovelChapterData) error
	ImportMetaData(novelName string, novelMetaData *NovelMetaData) error
	NumberOfChapter(novelName string) (int, error)

	ExportBook(novelName string, bookName string, content []byte) error
}
