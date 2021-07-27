package main

import "flag"

var (
	RSSListFilePath string
	OPMLFilePath    string
	RSS             string
	OutputFolder    string
	UserAgent       string
	ThreadCount     int
)

func init() {
	initFlags()
}

// initFlags defined all arguments used in PoDownloader
func initFlags() {
	flag.StringVar(&RSSListFilePath, "list", "", "Podcast RSS URL collection file path, one podcast RSS URL per line")
	flag.StringVar(&OPMLFilePath, "opml", "", "OPML file path")
	flag.StringVar(&RSS, "rss", "", "Podcast RSS URL")
	flag.StringVar(&OutputFolder, "o", "podcast", "Download destination folder")
	flag.StringVar(&UserAgent, "u", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36", "User agent")
	flag.IntVar(&ThreadCount, "t", 3, "Download threads")
}

// ParseFlags calls flag.Parse to parse all arguments
func ParseFlags() {
	flag.Parse()
}
