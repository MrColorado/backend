package models

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Chapter   int
	Paragraph []string
}

// TagData contain data on tag
type TagData struct {
	Id   int
	Name string
}

// BookData contain information on chapters' regroupment
type BookData struct {
	Id    int
	Start int
	End   int
}

type PartialNovelData struct {
	Id        int
	Title     string
	CoverPath string
	Tags      []TagData
}

type NovelData struct {
	CoreData        PartialNovelData
	NbChapter       int
	CurrentChapter  int
	Author          string
	FirstChapterURL string
	NextURL         string
	Summary         []string
}

// SCRAPER //

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
	Tags            []TagData
}
