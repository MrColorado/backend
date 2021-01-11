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
		Run: func(cmd *cobra.Command, args []string) {
			rootFunc()
		},
	}
)

func init() {
	rootCmd.AddCommand(partialCmd)

	partialCmd.Flags().IntVar(&partialOpt.startChapter, "start", 1, "Start chapter")
	partialCmd.Flags().IntVar(&partialOpt.endChapter, "end", 1, "End chapter")
}
