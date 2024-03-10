//go:generate sqlboiler -c ./../../../schemas/sqlboiler.toml --wipe psql

package dataStore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/dataStore/gen_models"
	"github.com/MrColorado/backend/book-handler/internal/models"
	"github.com/MrColorado/backend/logger"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type novelAuthor struct {
	gen_models.Novel `boil:",bind"`
	AuthorName       string `boil:"name"`
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(config config.PostgresConfigStruct) *PostgresClient {
	url := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", config.PostgresDB, config.PostgresUser, config.PostgresPassword, config.PostgresHost)
	db, err := sql.Open("postgres", url)
	if err != nil {
		logger.Warnf("failed to create connection to %s : %s", url, err.Error())
	}

	boil.SetDB(db)

	return &PostgresClient{
		db: db,
	}
}

// Todo update NovelMetaData to NovelData
func (client *PostgresClient) InsertOrUpdateNovel(data models.NovelMetaData) error {
	author := gen_models.Author{
		Name: data.Author,
	}
	err := author.Upsert(context.TODO(), client.db, false, nil, boil.Whitelist(), boil.Infer())
	if err != nil {
		return logger.Errorf("failed to upsert metadata for novel %s", data.Title)
	}
	if author.ID == 0 {
		atr, err := gen_models.Authors(gen_models.AuthorWhere.Name.EQ(data.Author)).One(context.TODO(), client.db)
		if err != nil {
			return logger.Errorf("failed to get author %s", data.Author)
		}
		author.ID = atr.ID
	}

	// Todo handle tags and genres
	novel := gen_models.Novel{
		Title:          data.Title,
		Summary:        data.Summary,
		FirstURL:       data.FirstURL,
		NextURL:        data.NextURL,
		CoverPath:      data.CoverPath,
		NBChapter:      data.NbChapter,
		CurrentChapter: data.CurrentChapter,
		LastUpdate:     time.Now(),
		FKAuthorID:     author.ID,
	}
	err = novel.Upsert(context.TODO(), client.db, true, []string{"title"}, boil.Greylist("CurrentChapter", "NextURL"), boil.Infer())
	if err != nil {
		return logger.Errorf("failed to upsert metadata for novel %s", data.Title)
	}
	return nil
}

func (client *PostgresClient) InsertOrUpdateBook(data models.BookData) error {
	// TODO migrate this to request in on sequence in order to keep data coherent even if one fail
	now := time.Now()
	book := gen_models.Book{
		FKNovelID:  data.NovelId,
		End:        data.End,
		Start:      data.Start,
		LastUpdate: now,
	}

	err := book.Upsert(context.TODO(), client.db, true, []string{"fk_novel_id", "start"}, boil.Greylist("end", "last_update"), boil.Infer())
	if err != nil {
		return logger.Errorf("failed to upsert book for novel %s starting %d and ending %d", data.NovelId, data.Start, data.End)
	}

	novel := gen_models.Novel{
		ID:         data.NovelId,
		LastUpdate: now,
	}
	_, err = novel.Update(context.TODO(), client.db, boil.Whitelist("last_update"))
	if err != nil {
		return logger.Errorf("failed to update last_update for for novel %s", data.NovelId)
	}

	return nil
}

func (client *PostgresClient) GetNovelByTitle(title string) (models.NovelData, error) {
	var na novelAuthor

	err := gen_models.NewQuery(
		qm.Select("author.name", "novel.id", "novel.nb_chapter", "novel.title", "novel.cover_path", "novel.first_url", "novel.next_url", "novel.current_chapter", "novel.summary", "novel.last_update"),
		qm.From("novel"),
		qm.InnerJoin("author on author.id = novel.fk_author_id"),
		qm.Where("novel.title = ?", title),
	).Bind(context.TODO(), client.db, &na)

	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel with title %s : %s", title, err.Error())
	}

	return models.NovelData{
		CoreData: models.PartialNovelData{
			Id:         na.ID,
			Title:      na.Title,
			Author:     na.AuthorName,
			CoverPath:  na.CoverPath,
			Summary:    na.Summary,
			LastUpdate: na.LastUpdate,
		},
		NbChapter:      na.NBChapter,
		CurrentChapter: na.CurrentChapter,
		FirstURL:       na.FirstURL,
		NextURL:        na.NextURL,
	}, nil
}
