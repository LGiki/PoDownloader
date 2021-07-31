package main

import (
	podownloader "PoDownloader"
	"PoDownloader/opml"
	"PoDownloader/podcast"
	"PoDownloader/util"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	HTTPClient    *http.Client
	PodcastParser *podcast.Parser
)

func init() {
	ParseFlags()
	HTTPClient = util.NewHTTPClient(UserAgent)
	PodcastParser = podcast.NewPodcastParser(HTTPClient)
}

// GetPodcastRSSList returns podcast rss URLs list parsed from opml file, RSS list file or RSS argument
func GetPodcastRSSList() ([]string, error) {
	if OPMLFilePath != "" {
		log.Println(fmt.Sprintf("Load podcast RSS links from OPML file: %s", OPMLFilePath))
		var (
			podcastOPML *opml.OPML
			err         error
		)
		if util.IsValidHttpLink(OPMLFilePath) {
			podcastOPML, err = opml.ParseOPMLFromURL(HTTPClient, OPMLFilePath)
		} else {
			podcastOPML, err = opml.ParseOPMLFromFile(OPMLFilePath)
		}
		if err != nil {
			return nil, err
		}
		return podcastOPML.GetAllXMLUrl(), nil
	} else if RSSListFilePath != "" {
		log.Println(fmt.Sprintf("Load podcast RSS links from RSS list file: %s", RSSListFilePath))
		return util.GetRSSListByTextFile(RSSListFilePath)
	} else if RSS != "" {
		log.Println(fmt.Sprintf("Load podcast RSS link from RSS: %s", RSS))
		return []string{RSS}, nil
	}
	return nil, errors.New("")
}

func main() {
	if OPMLFilePath == "" && RSSListFilePath == "" && RSS == "" {
		log.Fatalln("Please specify at least one parameter among opml, l and r")
	}
	podcastRSSList, err := GetPodcastRSSList()
	if err != nil {
		log.Fatalln("Can not load RSS list:", err)
	}
	log.Println(fmt.Sprintf("Loaded %d RSS link(s)", len(podcastRSSList)))

	podcastList, failed := PodcastParser.ParsePodcastsFromRSSList(podcastRSSList)
	log.Println(fmt.Sprintf("Successfully parsed %d RSS link(s)", len(podcastList)))
	if len(failed) != 0 {
		log.Println(fmt.Sprintf("%d RSS link(s) parsing failed:", len(failed)))
		for index, rss := range failed {
			log.Println(fmt.Sprintf("%d. %s", index+1, rss))
		}
	}

	// Exit when there are no podcasts to download
	if len(podcastList) == 0 {
		log.Println("No RSS to download, exit")
		os.Exit(0)
	}

	var podcastDownloadTasks []*podownloader.PodcastDownloadTask
	for _, p := range podcastList {
		tasks := p.GetPodcastDownloadTask(OutputFolder, HTTPClient)
		podcastDownloadTasks = append(podcastDownloadTasks, tasks)
	}
	podcastDownloadTaskIterator := podownloader.NewDownloadTaskIterator(podcastDownloadTasks)
	podcastDownloadTaskIterator.RemoveDownloadedTaskAndMakeDir(ThreadCount)

	log.Println("Start download")
	downloadQueue := podownloader.NewDownloadQueueFromDownloadTasks(podcastDownloadTaskIterator.PodcastDownloadTasks)
	failedTaskDestPaths := downloadQueue.StartDownload(HTTPClient, ThreadCount)
	log.Println("Download finished")

	// Print failed tasks
	if len(failedTaskDestPaths) > 0 {
		fmt.Println(fmt.Sprintf("%d file(s) download failed:", len(failedTaskDestPaths)))
		for index, failedTaskDestPath := range failedTaskDestPaths {
			fmt.Println(fmt.Sprintf("%d. %s", index+1, failedTaskDestPath))
		}
	}
}
