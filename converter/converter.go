package converter

// Converter interface that convert novelData to e-tablet format
type Converter interface {
	ConvertNovel(novelName string)
	ConvertPartialNovel(novelName string, startChapter int, endChapter int)
}
