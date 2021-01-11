package cmd

import (
	"github.com/spf13/cobra"
)

type PartialOpt struct {
	startChapter int
	endChapter   int
}

var (
	partialOpt = PartialOpt{}

	partialCmd = &cobra.Command{
		Use:   "partial",
		Short: "Indicate between wich chapter action must be done",
	}
)

func init() {
	rootCmd.AddCommand(partialCmd)

	partialCmd.Flags().IntVar(&partialOpt.startChapter, "start", 0, "Start chapter")
	partialCmd.Flags().IntVar(&partialOpt.endChapter, "end", 0, "End chapter")
}
