package models

import "time"

// NovelChapterData contain data of a chapter
type NovelChapterData struct {
	Chapter   int
	Paragraph []string
}

// BookData contain information on chapters' regroupment
type BookData struct {
	Id      string
	NovelId string
	Start   int
	End     int
}

type PartialNovelData struct {
	Id         string
	Title      string
	Author     string
	CoverPath  string
	Summary    string
	Genres     []string
	LastUpdate time.Time
}

type NovelData struct {
	CoreData       PartialNovelData
	NbChapter      int
	CurrentChapter int
	FirstURL       string
	NextURL        string
	Tags           []string
}
