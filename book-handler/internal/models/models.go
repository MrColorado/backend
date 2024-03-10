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

// SCRAPER //

// NovelMetaData contain data on the novel
type NovelMetaData struct {
	Title          string
	Author         string
	CoverPath      string
	FirstURL       string
	Summary        string
	NextURL        string
	NbChapter      int
	CurrentChapter int
	CoverData      []byte
	Genres         []string
	Tags           []string
	LastUpdate     time.Time
}

func MetaToNovel(data NovelMetaData) NovelData {
	return NovelData{
		CoreData: PartialNovelData{
			Title:      data.Title,
			Author:     data.Author,
			CoverPath:  data.CoverPath,
			Summary:    data.Summary,
			Genres:     data.Genres,
			LastUpdate: data.LastUpdate,
		},
		NbChapter:      data.NbChapter,
		CurrentChapter: data.CurrentChapter,
		FirstURL:       data.FirstURL,
		NextURL:        data.NextURL,
	}
}

func NovelToMeta(data NovelData) NovelMetaData {
	return NovelMetaData{
		Title:          data.CoreData.Title,
		Author:         data.CoreData.Author,
		CoverPath:      data.CoreData.CoverPath,
		Summary:        data.CoreData.Summary,
		Genres:         data.CoreData.Genres,
		FirstURL:       data.FirstURL,
		NextURL:        data.NextURL,
		NbChapter:      data.NbChapter,
		CurrentChapter: data.CurrentChapter,
		Tags:           data.Tags,
		CoverData:      []byte{},
	}
}
