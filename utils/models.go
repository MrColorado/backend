package utils

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Title     string
	Chapter   int
	Paragraph []string
}

// NovelMetaData contain data on the novel
type NovelMetaData struct {
	Title           string
	Author          string
	Summary         []string
	NumberOfChapter int
	ImagePath       string
}
