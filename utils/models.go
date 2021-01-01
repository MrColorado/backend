package utils

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Title     string
	Chapter   int
	Paragraph []string
}

// NovelData contain data for a (part of) novel
type NovelData struct {
	Author          string
	NumberOfChapter int
	Summary         []string
	Chapters        []NovelChapterData
}
