package models

import "time"

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Chapter   int
	Paragraph []string
}

// BookData contain information on chapters' regroupment
type BookData struct {
	NovelId int
	Start   int
	End     int
}

type PartialNovelData struct {
	Title      string
	Author     string
	CoverPath  string
	Summary    string
	Genres     []string
	LastUpdate time.Time
}

type NovelData struct {
	CoreData        PartialNovelData
	NbChapter       int
	CurrentChapter  int
	FirstChapterURL string
	NextURL         string
}

// SCRAPER //

// NovelMetaData contain data on the novel
type NovelMetaData struct {
	NbChapter       int
	Title           string
	Author          string
	CoverPath       string
	FirstChapterURL string
	NextURL         string
	CurrentChapter  int
	CoverData       []byte
	Summary         []string
	Genres          []string
	Tags            []string
}
