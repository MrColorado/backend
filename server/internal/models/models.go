package models

import "time"

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Paragraph []string
	Chapter   int
}

// BookData contain information on chapters' regroupment
type BookData struct {
	ID      string
	NovelID string
	Start   int
	End     int
}

type PartialNovelData struct {
	LastUpdate time.Time
	ID         string
	Title      string
	Author     string
	CoverPath  string
	Summary    string
	Genres     []string
}

type NovelData struct {
	CoreData       PartialNovelData
	FirstURL       string
	NextURL        string
	Tags           []string
	NbChapter      int
	CurrentChapter int
}
