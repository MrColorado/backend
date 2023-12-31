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

func (client *PostgresClient) InsertOrUpdateNovel(data models.NovelMetaData) error {
	novel := gen_models.Novel{
		Title:          data.Title,
		Author:         data.Author,
		Description:    strings.Join(data.Summary, "\n"),
		NBChapter:      data.NbChapter,
		FirstChapter:   data.FirstChapterURL,
		CurrentChapter: data.CurrentChapter,
		NextURL:        data.NextURL,
		LastUpdate:     time.Now(),
		CoverPath:      "",
	}
	err := novel.Upsert(context.TODO(), client.db, true, []string{"title", "author"}, boil.Greylist("CurrentChapter", "NextURL"), boil.Infer())
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

func (client *PostgresClient) GetNovelByTitle(title string) (models.NovelMetaData, error) {
	novel, err := gen_models.Novels(gen_models.NovelWhere.Title.EQ(title)).One(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return models.NovelMetaData{}, fmt.Errorf("failed to get novel with title %s", title)
	}

	return models.NovelMetaData{
		Id:              novel.ID,
		Title:           novel.Title,
		Author:          novel.Author,
		NbChapter:       novel.NBChapter,
		FirstChapterURL: novel.FirstChapter,
		Summary:         []string{novel.Description},
		CurrentChapter:  novel.CurrentChapter,
		NextURL:         novel.NextURL,
	}, nil
}

func (client *PostgresClient) GetNovelById(id int) (models.NovelMetaData, error) {
	novel, err := gen_models.Novels(gen_models.NovelWhere.ID.EQ(id)).One(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return models.NovelMetaData{}, fmt.Errorf("failed to get novel with id %d", id)
	}

	return models.NovelMetaData{
		Id:              novel.ID,
		Title:           novel.Title,
		Author:          novel.Author,
		NbChapter:       novel.NBChapter,
		FirstChapterURL: novel.FirstChapter,
		Summary:         []string{novel.Description},
		CurrentChapter:  novel.CurrentChapter,
		NextURL:         novel.NextURL,
	}, nil
}

func (client *PostgresClient) ListNovels() ([]models.PartialNovelData, error) {
	res := []models.PartialNovelData{}

	novels, err := gen_models.Novels().All(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return []models.PartialNovelData{}, fmt.Errorf("failed to get novels")
	}

	for _, novel := range novels {
		tagsDB, err := gen_models.Tags(
			qm.InnerJoin("novel_tag_map ntm on ntm.fk_tag_id = tag.id"),
			qm.Where("ntm.fk_novel_id=?", novel.ID),
		).All(context.TODO(), client.db)

		if err != nil {
			fmt.Println(err.Error())
			return []models.PartialNovelData{}, fmt.Errorf("failed to get tags for novel : %s", novel.Title)
		}

		tags := []models.TagData{}
		for _, tag := range tagsDB {
			tags = append(tags, models.TagData{Id: tag.ID, Name: tag.Name})
		}
		res = append(res, models.PartialNovelData{
			Id:        novel.ID,
			Title:     novel.Title,
			CoverPath: novel.CoverPath,
			Tags:      tags,
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
			Id:    book.ID,
			Start: book.Start,
			End:   book.End,
		})
	}
	return res, nil
}
