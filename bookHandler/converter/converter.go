package converter

// Converter interface that convert novelData to e-tablet format
type Converter interface {
	ConvertNovel(novelName string) error
	ConvertPartialNovel(novelName string, startChapter int, endChapter int) error
}
