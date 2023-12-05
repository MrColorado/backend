package models

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Chapter   int
	Paragraph []string
}

// NovelMetaData contain data on the novel
type NovelMetaData struct {
	Id              int
	NbChapter       int
	Title           string
	Author          string
	Summary         []string
	ImagePath       string
	FirstChapterURL string
	NextURL         string
	CurrentChapter  int
}
