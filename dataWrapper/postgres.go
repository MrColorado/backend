//go:generate sqlboiler -c schemas/sqlboiler.toml --wipe psql

package dataWrapper

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/dataWrapper/gen_models"
	"github.com/MrColorado/epubScraper/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type novelAndTag struct {
	Novel gen_models.Novel `boil:",bind"`
	Tag   gen_models.Tag   `boil:",bind"`
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

func (client *PostgresClient) InsertOrUpdateNovel(data models.NovelMetaData) error {
	novel := gen_models.Novel{
		Title:          data.Title,
		NBChapter:      data.NbChapter,
		FirstChapter:   data.FirstChapterURL,
		Author:         data.Author,
		Description:    strings.Join(data.Summary, "\n"),
		CurrentChapter: data.CurrentChapter,
		NextURL:        data.NextURL,
	}
	err := novel.Upsert(context.TODO(), client.db, true, []string{"title"}, boil.Greylist("CurrentChapter", "NextURL"), boil.Infer())
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to upsert metadata for novel %s", data.Title)
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

func partialNovelFromNovelAndTag(nat novelAndTag) models.PartialNovelData {
	return models.PartialNovelData{
		Id:        nat.Novel.ID,
		Title:     nat.Novel.Title,
		CoverPath: nat.Novel.CoverPath,
		Tags: []models.TagData{
			{
				Id:   nat.Tag.ID,
				Name: nat.Tag.Name,
			},
		},
	}
}

func (client *PostgresClient) ListNovels() ([]models.PartialNovelData, error) {
	var novelsAndTags []novelAndTag

	err := gen_models.NewQuery(
		qm.From("novel"),
		qm.Select("novel.id", "novel.title", "novel.cover_path", "tag.id", "tag.name"),
		qm.InnerJoin("novel_tag_map ntm on ntm.fk_novel_id = novel.id"),
		qm.InnerJoin("tag on tag.id = ntm.fk_tag_id"),
	).Bind(context.TODO(), client.db, &novelsAndTags)

	if err != nil {
		fmt.Println(err)
		return []models.PartialNovelData{}, fmt.Errorf("failed to list novels")
	} else if len(novelsAndTags) == 0 {
		return []models.PartialNovelData{}, nil
	}

	res := []models.PartialNovelData{}
	var novel models.PartialNovelData
	for pos, nat := range novelsAndTags {
		if pos == 0 {
			novel = partialNovelFromNovelAndTag(nat)
			continue
		}
		if novel.Id != nat.Novel.ID {
			res = append(res, novel)
			novel = partialNovelFromNovelAndTag(nat)
			continue
		}
		novel.Tags = append(novel.Tags, models.TagData{Id: nat.Tag.ID, Name: nat.Tag.Name})
	}
	res = append(res, novel)
	return res, nil
}

func (client *PostgresClient) ListBooks(novelId int) ([]models.BookData, error) {
	books, err := gen_models.Books(gen_models.BookWhere.ID.EQ(novelId)).All(context.TODO(), client.db)

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
