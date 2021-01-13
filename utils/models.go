package utils

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Chapter   int
	Paragraph []string
}

// NovelMetaData contain data on the novel
type NovelMetaData struct {
	Title           string
	Author          string
	Summary         []string
	ImagePath       string
	FirstChapterURL string
	NextURL         string
}
