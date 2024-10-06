//go:generate sqlboiler -c ./../../../schemas/sqlboiler.toml --wipe psql

package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/data/gen_models"
	"github.com/MrColorado/backend/book-handler/internal/models"
	"github.com/MrColorado/backend/logger"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type novelAuthor struct {
	AuthorName       string `boil:"name"`
	gen_models.Novel `boil:",bind"`
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(cfg config.PostgresConfigStruct) *PostgresClient {
	url := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", cfg.PostgresDB, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost)
	db, err := sql.Open("postgres", url)
	if err != nil {
		logger.Warnf("failed to create connection to %s : %s", url, err.Error())
	}

	boil.SetDB(db)

	return &PostgresClient{
		db: db,
	}
}

func (client *PostgresClient) InsertOrUpdateNovel(data models.NovelMetaData, genre bool) error {
	// Todo update NovelMetaData to NovelData
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

	if genre {
		genres := []*gen_models.Genre{}
		for _, genre := range data.Genres {
			genres = append(genres, &gen_models.Genre{Name: genre})
		}

		err = novel.AddFKGenreNameGenres(context.TODO(), client.db, false, genres...)
		if err != nil {
			logger.Fatalf("failed to add genre for novel %s : %s", data.Title, err.Error())
		}
	}

	return nil
}

func (client *PostgresClient) InsertOrUpdateBook(data models.BookData) error {
	// TODO migrate this to request in on sequence in order to keep data coherent even if one fail
	now := time.Now()
	book := gen_models.Book{
		FKNovelID:  data.NovelID,
		End:        data.End,
		Start:      data.Start,
		LastUpdate: now,
	}

	err := book.Upsert(context.TODO(), client.db, true, []string{"fk_novel_id", "start"}, boil.Greylist("end", "last_update"), boil.Infer())
	if err != nil {
		return logger.Errorf("failed to upsert book for novel %s starting %d and ending %d", data.NovelID, data.Start, data.End)
	}

	novel := gen_models.Novel{
		ID:         data.NovelID,
		LastUpdate: now,
	}
	_, err = novel.Update(context.TODO(), client.db, boil.Whitelist("last_update"))
	if err != nil {
		return logger.Errorf("failed to update last_update for novel %s", data.NovelID)
	}

	return nil
}

func (client *PostgresClient) InsertOrUpdateGenre(name string) error {
	genre := gen_models.Genre{
		Name: name,
	}

	err := genre.Upsert(context.TODO(), client.db, false, []string{}, boil.Infer(), boil.Infer())

	if err != nil {
		return logger.Errorf("failed to add genre %s : %s", name, err.Error())
	}

	return nil
}

func (client *PostgresClient) GetNovelByTitle(title string) (models.NovelData, error) {
	var na novelAuthor

	err := gen_models.NewQuery(
		qm.Select("author.name as name", "novel.id", "novel.nb_chapter", "novel.title", "novel.cover_path", "novel.first_url", "novel.next_url", "novel.current_chapter", "novel.summary", "novel.last_update"), //nolint:lll
		qm.From("novel"),
		qm.InnerJoin("author on author.id = novel.fk_author_id"),
		qm.Where("novel.title = ?", title),
	).Bind(context.TODO(), client.db, &na)

	if err != nil {
		logger.Warnf("failed to get novel with title %s : %s", title, err.Error())
		return models.NovelData{}, nil
	}

	return models.NovelData{
		CoreData: models.PartialNovelData{
			ID:         na.ID,
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

func (client *PostgresClient) GetBookByTitle(title string, start int) (models.BookData, error) {
	var book gen_models.Book
	err := gen_models.NewQuery(
		qm.Select("novel.id", "fk_novel_id", "start", "end"),
		qm.From("book"),
		qm.InnerJoin("novel on novel.id = book.fk_novel_id"),
		qm.Where("novel.title = ?", title),
		qm.Where("start = ?", start),
	).Bind(context.TODO(), client.db, &book)

	if err != nil {
		return models.BookData{}, err
	}

	return models.BookData{
		NovelID: book.FKNovelID,
		Start:   book.Start,
		End:     book.End,
	}, nil
}
