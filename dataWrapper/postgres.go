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
)

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(config configuration.PostgresConfigStruct) *PostgresClient {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", config.PostgresDB, config.PostgresUser, config.PostgresPassword))
	if err != nil {
		fmt.Println(err)
	}

	boil.SetDB(db)

	return &PostgresClient{
		db: db,
	}
}

func (client *PostgresClient) InsertOrUpdate(data *models.NovelMetaData) error {
	novel := gen_models.Novel{
		ID:             data.ID,
		Title:          data.Title,
		NBChapter:      data.NbChapter,
		FirstChapter:   data.FirstChapterURL,
		Author:         data.Author,
		Description:    strings.Join(data.Summary, "\n"),
		CurrentChapter: data.CurrentChapter,
		NextURL:        data.NextURL,
	}
	err := novel.Upsert(context.TODO(), client.db, true, []string{"id"}, boil.Greylist("CurrentChapter", "NextURL"), boil.Infer())
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to upsert metadata for novel %s", data.Title)
	}
	data.ID = novel.ID
	return nil
}

func (client *PostgresClient) GetId(id int) (models.NovelMetaData, error) {
	novel, err := gen_models.Novels(gen_models.NovelWhere.ID.EQ(id)).One(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return models.NovelMetaData{}, fmt.Errorf("failed to get novel with ID %d", id)
	}

	return models.NovelMetaData{
		ID:              novel.ID,
		Title:           novel.Title,
		Author:          novel.Author,
		NbChapter:       novel.NBChapter,
		FirstChapterURL: novel.FirstChapter,
		Summary:         []string{novel.Description},
		CurrentChapter:  novel.CurrentChapter,
		NextURL:         novel.NextURL,
	}, nil
}

func (client *PostgresClient) GetTitle(title string) (models.NovelMetaData, error) {
	novel, err := gen_models.Novels(gen_models.NovelWhere.Title.EQ(title)).One(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return models.NovelMetaData{}, fmt.Errorf("failed to get novel with title %s", title)
	}

	return models.NovelMetaData{
		ID:              novel.ID,
		Title:           novel.Title,
		Author:          novel.Author,
		NbChapter:       novel.NBChapter,
		FirstChapterURL: novel.FirstChapter,
		Summary:         []string{novel.Description},
		CurrentChapter:  novel.CurrentChapter,
		NextURL:         novel.NextURL,
	}, nil
}

func (client *PostgresClient) List() ([]models.NovelMetaData, error) {
	novels, err := gen_models.Novels().All(context.TODO(), client.db)
	if err != nil {
		fmt.Println(err)
		return []models.NovelMetaData{}, fmt.Errorf("failed to list novels")
	}

	res := []models.NovelMetaData{}
	for _, novel := range novels {
		res = append(res, models.NovelMetaData{
			ID:              novel.ID,
			Title:           novel.Title,
			Author:          novel.Author,
			NbChapter:       novel.NBChapter,
			FirstChapterURL: novel.FirstChapter,
			Summary:         []string{novel.Description},
			CurrentChapter:  novel.CurrentChapter,
			NextURL:         novel.NextURL,
		})
	}
	return res, nil
}
