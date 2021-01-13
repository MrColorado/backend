package cmd

import (
	"fmt"
	"os"

	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

type RootOpt struct {
	action      Flags
	websiteName Flags
	novelName   string
	outputPath  string
}

var (
	websiteNames = map[int32]string{
		0: "READ_NOVEL_FULL",
	}
	websiteValues = map[string]int32{
		"READ_NOVEL_FULL": 0,
	}
	actionNames = map[int32]string{
		0: "GENERATE",
		1: "SCRAPE",
		2: "CONVERT",
	}
	actionValues = map[string]int32{
		"GENERATE": 0,
		"SCRAPE":   1,
		"CONVERT":  2,
	}
)

var (
	rootOpt = RootOpt{
		action: Flags{
			EnumName:  actionNames,
			EnumValue: actionValues,
			Usage:     "Action to produce",
		},
		websiteName: Flags{
			EnumName:  websiteNames,
			EnumValue: websiteValues,
			Usage:     "Name of the website",
		},
	}

	rootCmd = &cobra.Command{
		Short:            "A scraper of novel's aggregator website",
		Long:             `NovelRecuperator is a software that allow you to get any novel avaible on some websites in EPUB format`,
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			rootFunc()
		},
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Var(&rootOpt.action, "action", rootOpt.action.Usage)
	rootCmd.Flags().Var(&rootOpt.websiteName, "website-name", rootOpt.websiteName.Usage)
	rootCmd.Flags().StringVar(&rootOpt.novelName, "novel-name", "", "Name of the novel")
	rootCmd.Flags().StringVar(&rootOpt.outputPath, "output-path", "~/Novels", "Location of files")
}

func rootFunc() {
	if actionValues["SCRAPE"] == rootOpt.action.Value {
		scrape()
	} else if actionValues["CONVERT"] == rootOpt.action.Value {
		convert()
	} else {
		generate()
	}
}

func scrape() {
	c := colly.NewCollector()
	scraper := scraper.ReadNovelScraper{}
	if partialOpt.endChapter > 1 {
		scraper.ScrapPartialNovel(c, rootOpt.novelName, fmt.Sprintf("%s/raw", rootOpt.outputPath),
			partialOpt.startChapter, partialOpt.endChapter)
	} else {
		scraper.ScrapeNovel(c, rootOpt.novelName, fmt.Sprintf("%s/raw", rootOpt.outputPath))
	}
}

func convert() {
	converter := converter.EpubConverter{}
	if partialOpt.endChapter > 1 {
		converter.ConvertPartialNovel(fmt.Sprintf("%s/raw", rootOpt.outputPath), fmt.Sprintf("%s/epub", rootOpt.outputPath),
			rootOpt.novelName, partialOpt.startChapter, partialOpt.endChapter)
	} else {
		converter.ConvertNovel(fmt.Sprintf("%s/raw", rootOpt.outputPath),
			fmt.Sprintf("%s/epub", rootOpt.outputPath), rootOpt.novelName)
	}
}

func generate() {
	scrape()
	convert()
}
