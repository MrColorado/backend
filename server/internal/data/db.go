//go:generate sqlboiler -c ./../../../schemas/sqlboiler.toml --wipe psql

package data

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/MrColorado/backend/logger"
	"github.com/MrColorado/backend/server/internal/config"
	"github.com/MrColorado/backend/server/internal/data/gen_models"
	"github.com/MrColorado/backend/server/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type novelAuthor struct {
	AuthorName       string `boil:"name"`
	gen_models.Novel `boil:",bind"`
}

type novelGenres struct {
	GenreName        string `boil:"genre_name"`
	AuthorName       string `boil:"author_name"`
	gen_models.Novel `boil:",bind"`
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(cfg config.PostgresConfigStruct) *PostgresClient {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", cfg.PostgresDB, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost))
	if err != nil {
		logger.Warn(err.Error())
	}

	boil.SetDB(db)
	return &PostgresClient{
		db: db,
	}
}

func (client *PostgresClient) GetNovelByID(id string) (models.NovelData, error) {
	var na novelAuthor

	err := gen_models.NewQuery(
		qm.Select("author.name as name", "novel.id", "novel.nb_chapter", "novel.title", "novel.cover_path", "novel.first_url", "novel.next_url", "novel.current_chapter", "novel.summary", "novel.last_update"), //nolint:lll
		qm.From("novel"),
		qm.InnerJoin("author on author.id = novel.fk_author_id"),
		qm.Where("novel.id = ?", id),
	).Bind(context.TODO(), client.db, &na)

	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel with id %s", id)
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

func (client *PostgresClient) GetNovelByTitle(title string) (models.NovelData, error) {
	var na novelAuthor

	err := gen_models.NewQuery(
		qm.Select("author.name", "novel.id", "novel.nb_chapter", "novel.title", "novel.cover_path", "novel.first_url", "novel.next_url", "novel.current_chapter", "novel.summary", "novel.last_update"),
		qm.From("novel"),
		qm.InnerJoin("author on author.id = novel.fk_author_id"),
		qm.Where("novel.title = ?", title),
	).Bind(context.TODO(), client.db, &na)

	if err != nil {
		return models.NovelData{}, logger.Errorf("failed to get novel with title %s", title)
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

func (client *PostgresClient) ListNovels(title string) ([]models.PartialNovelData, error) {
	res := []models.PartialNovelData{}

	var ng []novelGenres
	err := gen_models.NewQuery(
		qm.Select("ngm.fk_genre_name as genre_name", "author.name as author_name", "novel.id", "novel.title", "novel.cover_path", "novel.summary", "novel.last_update"),
		qm.From("novel"),
		qm.InnerJoin("novel_genre_map ngm on ngm.fk_novel_id = novel.id"),
		qm.InnerJoin("author on author.id = novel.fk_author_id"),
		qm.Where("title like ?", fmt.Sprintf("%%%s%%", title)),
	).Bind(context.TODO(), client.db, &ng)

	if err != nil {
		return []models.PartialNovelData{}, logger.Errorf("failed to get novels")
	}

	for i := 0; i < len(ng); {
		pnd := models.PartialNovelData{
			ID:         ng[i].ID,
			Title:      ng[i].Title,
			CoverPath:  ng[i].CoverPath,
			Summary:    ng[i].Summary,
			LastUpdate: ng[i].LastUpdate,
			Author:     ng[i].AuthorName,
		}

		genres := []string{}
		for ; i < len(ng) && ng[i].ID == pnd.ID; i++ {
			genres = append(genres, ng[i].GenreName)
		}
		pnd.Genres = genres

		res = append(res, pnd)
	}

	return res, nil
}

func (client *PostgresClient) ListBooks(novelID string) ([]models.BookData, error) {
	books, err := gen_models.Books(gen_models.BookWhere.FKNovelID.EQ(novelID)).All(context.TODO(), client.db)

	if err != nil {
		return []models.BookData{}, logger.Errorf("failed to list novels")
	}

	res := []models.BookData{}
	for _, book := range books {
		res = append(res, models.BookData{
			Start: book.Start,
			End:   book.End,
		})
	}
	return res, nil
}
