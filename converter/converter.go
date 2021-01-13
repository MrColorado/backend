package converter

// Converter interface that convert novelData to e-tablet format
type Converter interface {
	ConvertNovel(inputPath string, outputPath string, novelName string)
	ConvertPartialNovel(inputPath string, outputPath string, novelName string, startChapter int, endChapter int)
}
