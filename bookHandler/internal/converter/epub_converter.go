package converter

import (
	"bytes"
	"fmt"

	"github.com/MrColorado/backend/bookHandler/internal/core"
	"github.com/MrColorado/backend/bookHandler/internal/models"
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

func (cvt *EpubConverter) convertMetaData(e *epub.Epub, novelName string) error {
	data, err := cvt.app.GetMetaData(novelName)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("failed to import metaData for novel %s", novelName)
	}

	filePath, err := cvt.app.GetCoverDiskPath(novelName)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("failed to import cover inside novel %s", novelName)
	}
	// defer cvt.app.RemoveCoverDiskPath(filePath) TODO check if we can delete the image before the novel is written on disk

	imgPath, err := e.AddImage(filePath, "")
	if err != nil {
		println(err.Error())
		return fmt.Errorf("failed to import cover inside novel %s", novelName)
	}
	coverCSSPath, err := e.AddCSS("./converter/epub.css", "")
	if err != nil {
		println(err.Error())
		return fmt.Errorf("failed to import cover inside novel %s", novelName)
	}

	e.SetAuthor(data.Author)
	e.SetTitle(novelName)
	e.SetCover(imgPath, coverCSSPath)
	e.SetDescription(fmt.Sprintf("<p>%s</p>", data.Summary))
	return nil
}

func (cvt *EpubConverter) convertToNovel(novelName string, startChapter int, endChapter int) error {
	fileName := fmt.Sprintf("%s-%04d-%04d", novelName, startChapter, endChapter)
	e := epub.NewEpub(fileName)
	cvt.convertMetaData(e, novelName)

	for i := startChapter; i <= endChapter; i++ {
		chapterData, err := cvt.app.GetNovelChapter(novelName, i)
		chapterData.Chapter = i
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		bodySection := fmt.Sprintf("<h1>Chapter %d</h1>", i)
		for _, paragraph := range chapterData.Paragraph {
			bodySection += fmt.Sprintf("<p>%s</p>", paragraph)
		}
		if _, err := e.AddSection(bodySection, fmt.Sprintf("Chapter %d", i), "", ""); err != nil {
			fmt.Printf("fail to add chapter %d of novel %s\n", i, novelName)
			continue
		}
	}

	buf := new(bytes.Buffer)
	_, err := e.WriteTo(buf)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to write epub file %s in a buffer", fileName)
	}

	err = cvt.app.ExportBook(novelName, fileName, buf.Bytes(), models.BookData{Start: startChapter, End: endChapter})
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export epub file %s", fileName)
	}
	return nil
}

// ConvertPartialNovel convert partial novel to epub (startChapter include / endChapter included)
func (cvt *EpubConverter) ConvertPartialNovel(novelName string, startChapter int, endChapter int) error {
	if endChapter > 100 && startChapter%100 != 1 {
		toModulo100 := 100 - startChapter%100
		err := cvt.convertToNovel(novelName, startChapter, startChapter+toModulo100)
		if err != nil {
			fmt.Println(err.Error())
		}
		startChapter += toModulo100
	}

	numberOfBook := 1 + (endChapter-startChapter)/100
	firstBook := startChapter / 100

	for i := firstBook; i < firstBook+numberOfBook; i++ {
		err := cvt.convertToNovel(novelName, i*100+1, (i+1)*100)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if endChapter%100 != 0 {
		err := cvt.convertToNovel(novelName, (firstBook+numberOfBook)*100+1, endChapter)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}

// ConvertNovel convert every novel in inputPath to epub format
func (cvt *EpubConverter) ConvertNovel(novelName string) error {
	nbChapter, err := cvt.app.GetNbChapter(novelName)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to get number of chapter for novel %s", novelName)
	}

	rest := 0
	if nbChapter%100 != 0 {
		rest += 1
	}
	nbBook := nbChapter / 100
	fmt.Printf("for novel %s there are %d chapter and so %d books\n", novelName, nbChapter, nbBook+rest)

	for i := 0; i < nbBook; i++ {
		err := cvt.convertToNovel(novelName, (i*100)+1, (i+1)*100)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if rest != 0 {
		err := cvt.convertToNovel(novelName, (nbBook*100)+1, nbChapter)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}
