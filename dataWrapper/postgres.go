//go:generate sqlboiler -c schemas/sqlboiler.toml --wipe psql

package dataWrapper

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/dataWrapper/gen_models"
	"github.com/MrColorado/epubScraper/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type novelAuthor struct {
	gen_models.Novel  `boil:",bind"`
	gen_models.Author `boil:",bind"`
	// author           `boil:",author_name"`
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(config configuration.PostgresConfigStruct) *PostgresClient {
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", config.PostgresDB, config.PostgresUser, config.PostgresPassword, config.PostgresHost)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return fmt.Errorf("failed to upsert metadata for novel %s", data.Title)
	}

	// Todo handle tags and genres
	novel := gen_models.Novel{
		Title:          data.Title,
		FKAuthorID:     author.ID,
		Description:    strings.Join(data.Summary, "\n"),
		NBChapter:      data.NbChapter,
		FirstChapter:   data.FirstChapterURL,
		CurrentChapter: data.CurrentChapter,
		NextURL:        data.NextURL,
		LastUpdate:     time.Now(),
		CoverPath:      "",
	}
	err = novel.Upsert(context.TODO(), client.db, true, []string{"title", "author"}, boil.Greylist("CurrentChapter", "NextURL"), boil.Infer())
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to upsert metadata for novel %s", data.Title)
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
		fmt.Println(err)
		return fmt.Errorf("failed to upsert book for novel %d starting %d and ending %d", data.NovelId, data.Start, data.End)
	}

	novel := gen_models.Novel{
		ID:         data.NovelId,
		LastUpdate: now,
	}
	_, err = novel.Update(context.TODO(), client.db, boil.Whitelist("last_update"))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to update last_update for for novel %d", data.NovelId)
	}

	return nil
}

func (client *PostgresClient) GetNovelByTitle(title string) (models.NovelData, error) {
	var na novelAuthor

	err := gen_models.NewQuery(
		qm.Select("novel.nb_chapter", "novel.title", "author.name", "novel.cover_path", "novel.first_chapter", "novel.next_url", "novel.current_chapter", "novel.description", "novel.last_update"),
		qm.From("novel"),
		qm.InnerJoin("author on author.id = novel.fk_author_id"),
		qm.Where("novel.title = ?", title),
	).Bind(context.TODO(), client.db, &na)

	if err != nil {
		fmt.Println(err)
		return models.NovelData{}, fmt.Errorf("failed to get novel with title %s", title)
	}

	return models.NovelData{
		CoreData: models.PartialNovelData{
			Title:      na.Title,
			Author:     na.Name,
			CoverPath:  na.CoverPath,
			Summary:    na.Description,
			LastUpdate: na.LastUpdate,
		},
		NbChapter:       na.NBChapter,
		CurrentChapter:  na.CurrentChapter,
		FirstChapterURL: na.FirstChapter,
		NextURL:         na.NextURL,
	}, nil
}

func (client *PostgresClient) ListNovels(novelName string) ([]models.PartialNovelData, error) {
	res := []models.PartialNovelData{}

	novels, err := gen_models.Novels(qm.Where("title like ?", fmt.Sprintf("%%%s%%", novelName))).All(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return []models.PartialNovelData{}, fmt.Errorf("failed to get novels")
	}

	for _, novel := range novels {
		genreDb, err := gen_models.Genres(
			qm.InnerJoin("novel_genre_map ngm on ngm.fk_genre_id = genre.id"),
			qm.Where("ngm.fk_novel_id=?", novel.ID),
		).All(context.TODO(), client.db)

		if err != nil {
			fmt.Println(err.Error())
			return []models.PartialNovelData{}, fmt.Errorf("failed to get genres for novel : %s", novel.Title)
		}

		genres := []string{}
		for _, genre := range genreDb {
			genres = append(genres, genre.Name)
		}
		res = append(res, models.PartialNovelData{
			Title:      novel.Title,
			CoverPath:  novel.CoverPath,
			Summary:    novel.Description, // Todo description => summary
			Genres:     genres,
			LastUpdate: novel.LastUpdate,
		})
	}

	return res, nil
}

func (client *PostgresClient) ListBooks(novelId int) ([]models.BookData, error) {
	books, err := gen_models.Books(gen_models.BookWhere.FKNovelID.EQ(novelId)).All(context.TODO(), client.db)

	if err != nil {
		fmt.Println(err)
		return []models.BookData{}, fmt.Errorf("failed to list novels")
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
