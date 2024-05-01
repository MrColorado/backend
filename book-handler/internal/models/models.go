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
	ID         string
	Title      string
	Author     string
	CoverPath  string
	Summary    string
	LastUpdate time.Time
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

// SCRAPER //

// NovelMetaData contain data on the novel
type NovelMetaData struct {
	LastUpdate     time.Time
	Title          string
	Author         string
	CoverPath      string
	FirstURL       string
	Summary        string
	NextURL        string
	CoverData      []byte
	Genres         []string
	Tags           []string
	NbChapter      int
	CurrentChapter int
	Status         int
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
