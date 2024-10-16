package converter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/core"
	"github.com/MrColorado/backend/book-handler/internal/models"
	"github.com/MrColorado/backend/logger"
	"github.com/bmaupin/go-epub"
)

const (
	EpubConverterName = "Epub"
)

// EpubConverter convert novel data ton epub format
type EpubConverter struct {
	app *core.App
}

func NewEpubConverter(app *core.App) *EpubConverter {
	return &EpubConverter{
		app: app,
	}
}

func (cvt *EpubConverter) convertMetaData(e *epub.Epub, novelName string, bookName string) error {
	data, err := cvt.app.GetMetaData(novelName)
	if err != nil {
		return logger.Errorf("failed to import metaData for novel %s : %s", novelName, err.Error())
	}

	filePath, err := cvt.app.GetCoverDiskPath(novelName)
	if err != nil {
		return logger.Errorf("failed to import cover inside novel %s : %s", novelName, err.Error())
	}
	// defer cvt.app.RemoveCoverDiskPath(filePath) TODO check if we can delete the image before the novel is written on disk

	imgPath, err := e.AddImage(filePath, "")
	if err != nil {
		return logger.Errorf("failed to import cover inside novel %s : %s", novelName, err.Error())
	}

	coverCSSPath, err := e.AddCSS(config.GetConfig().MiscConfig.FilesFolder+"converter/css/epub.css", "")
	if err != nil {
		return logger.Errorf("failed to import cover inside novel %s : %s", novelName, err.Error())
	}

	e.SetAuthor(data.Author)
	e.SetTitle(bookName)
	e.SetCover(imgPath, coverCSSPath)
	e.SetDescription(fmt.Sprintf("<p>%s</p>", data.Summary))
	return nil
}

func (cvt *EpubConverter) convertToNovel(novelName string, startChapter int, endChapter int) error {
	fileName := fmt.Sprintf("%s-%04d-%04d", novelName, startChapter, endChapter)
	e := epub.NewEpub(fileName)
	cvt.convertMetaData(e, novelName, fmt.Sprintf("%s %04d-%04d", novelName, startChapter, endChapter))

	for i := startChapter; i <= endChapter; i++ {
		chapterData, err := cvt.app.GetNovelChapter(novelName, i)
		if err != nil {
			logger.Errorf("failed to get chapter %d of novel %s : %s", i, novelName, err.Error())
			continue
		}

		bodySection := fmt.Sprintf("<h1>Chapter %d</h1>", i)
		for _, paragraph := range chapterData.Paragraph {
			paragraph = strings.ReplaceAll(paragraph, "<", "--")
			paragraph = strings.ReplaceAll(paragraph, ">", "--")
			bodySection += fmt.Sprintf("<p>%s</p>", paragraph)
		}
		if _, err := e.AddSection(bodySection, fmt.Sprintf("Chapter %d", i), "", ""); err != nil {
			logger.Warnf("failed to add chapter %d of novel %s : %s", i, novelName, err.Error())
			continue
		}
	}

	buf := new(bytes.Buffer)
	_, err := e.WriteTo(buf)
	if err != nil {
		return logger.Errorf("failed to write epub file %s in a buffer", fileName)
	}

	err = cvt.app.ExportBook(novelName, fileName, buf.Bytes(), models.BookData{Start: startChapter, End: endChapter})
	if err != nil {
		return logger.Errorf("failed to export epub file %s", fileName)
	}
	return nil
}

// ConvertPartialNovel convert partial novel to epub (startChapter include / endChapter included)
func (cvt *EpubConverter) ConvertPartialNovel(novelName string, startChapter int, endChapter int) error {
	if endChapter > 100 && startChapter%100 != 1 {
		toModulo100 := 100 - startChapter%100
		err := cvt.convertToNovel(novelName, startChapter, startChapter+toModulo100)
		if err != nil {
			return logger.Errorf("failed to convert to novel %s, %d, %d", novelName, startChapter, startChapter+toModulo100)
		}
		startChapter += toModulo100
	}

	numberOfBook := 1 + (endChapter-startChapter)/100
	firstBook := startChapter / 100

	for i := firstBook; i < firstBook+numberOfBook; i++ {
		err := cvt.convertToNovel(novelName, i*100+1, (i+1)*100)
		if err != nil {
			return logger.Errorf("failed to convert to novel %s, %d, %d", novelName, i*100+1, (i+1)*100)
		}
	}

	if endChapter%100 != 0 {
		err := cvt.convertToNovel(novelName, (firstBook+numberOfBook)*100+1, endChapter)
		if err != nil {
			return logger.Errorf("failed to convert to novel %s, %d, %d", novelName, (firstBook+numberOfBook)*100+1, endChapter)
		}
	}
	return nil
}

// ConvertNovel convert every novel in inputPath to epub format
func (cvt *EpubConverter) ConvertNovel(novelName string) error {
	nbChapter, err := cvt.app.GetNbChapter(novelName)
	if err != nil {
		return logger.Errorf("failed to get number of chapter for novel %s", novelName)
	}

	rest := 0
	if nbChapter%100 != 0 {
		rest += 1
	}
	nbBook := nbChapter / 100
	logger.Infof("for novel %s there are %d chapter and so %d books", novelName, nbChapter, nbBook+rest)

	for i := 0; i < nbBook; i++ {
		err := cvt.convertToNovel(novelName, (i*100)+1, (i+1)*100)
		if err != nil {
			logger.Errorf("failed to convert novel %s, %d, %d", novelName, (i*100)+1, (i+1)*100)
		}
	}

	if rest != 0 {
		err := cvt.convertToNovel(novelName, (nbBook*100)+1, nbChapter)
		if err != nil {
			logger.Errorf("failed to convert novel %s, %d, %d", novelName, (nbBook*100)+1, nbChapter)
		}
	}
	return nil
}
