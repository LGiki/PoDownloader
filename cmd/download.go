package main

import (
	podownloader "PoDownloader"
	"PoDownloader/opml"
	"PoDownloader/podcast"
	"PoDownloader/util"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

var (
	// httpClient and podcastParser are used to download podcasts
	httpClient    *http.Client
	podcastParser *podcast.Parser

	// arguments used in download command
	rssListFilePath string
	opmlFilePath    string
	rss             string
	outputFolder    string
	userAgent       string
	configFilePath  string
	threadCount     int

	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download podcasts",
		Run:   download,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	httpClient = util.NewHTTPClient(userAgent)
	podcastParser = podcast.NewPodcastParser(httpClient)
	rootCmd.AddCommand(downloadCmd)

	// Define download command flags
	downloadCmd.Flags().StringVarP(&rssListFilePath, "list", "l", "", "Podcast RSS URL collection file path, one podcast RSS URL per line")
	downloadCmd.Flags().StringVarP(&opmlFilePath, "opml", "f", "", "OPML file path")
	downloadCmd.Flags().StringVarP(&rss, "rss", "r", "", "Podcast RSS URL")
	downloadCmd.Flags().StringVarP(&outputFolder, "output", "o", "podcast", "Download destination folder")
	downloadCmd.Flags().StringVarP(&userAgent, "ua", "u", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36", "User Agent")
	downloadCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "Configuration file (default is $PWD/.podownloader)")
	downloadCmd.Flags().IntVarP(&threadCount, "thread", "t", 3, "Download threads")

	// Define configuration file keys
	_ = viper.BindPFlag("list", rootCmd.Flags().Lookup("list"))
	_ = viper.BindPFlag("opml", rootCmd.Flags().Lookup("opml"))
	_ = viper.BindPFlag("rss", rootCmd.Flags().Lookup("rss"))
	_ = viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("ua", rootCmd.Flags().Lookup("ua"))
	_ = viper.BindPFlag("thread", rootCmd.Flags().Lookup("thread"))
	viper.SetDefault("output", "podcast")
	viper.SetDefault("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	viper.SetDefault("thread", 3)
}

// getPodcastRSSList returns podcast rss URLs list parsed from OPML file, RSS list file or RSS argument
func getPodcastRSSList() ([]string, error) {
	if opmlFilePath != "" {
		log.Println(fmt.Sprintf("Load podcast RSS links from OPML file: %s", opmlFilePath))
		var (
			podcastOPML *opml.OPML
			err         error
		)
		if util.IsValidHttpLink(opmlFilePath) {
			podcastOPML, err = opml.ParseOPMLFromURL(httpClient, opmlFilePath)
		} else {
			podcastOPML, err = opml.ParseOPMLFromFile(opmlFilePath)
		}
		if err != nil {
			return nil, err
		}
		return podcastOPML.GetAllXMLUrl(), nil
	} else if rssListFilePath != "" {
		log.Println(fmt.Sprintf("Load podcast RSS links from RSS list file: %s", rssListFilePath))
		return util.GetRSSListByTextFile(rssListFilePath)
	} else if rss != "" {
		log.Println(fmt.Sprintf("Load podcast from RSS: %s", rss))
		return []string{rss}, nil
	}
	return nil, errors.New("")
}

func download(cmd *cobra.Command, _ []string) {
	if opmlFilePath == "" && rssListFilePath == "" && rss == "" {
		fmt.Println("Please specify at least one argument among \"opml\", \"list\" and \"rss\"")
		fmt.Println()
		_ = cmd.Help()
		os.Exit(1)
	}
	podcastRSSList, err := getPodcastRSSList()
	if err != nil {
		log.Fatalln("Can not load RSS list:", err)
	}
	podcastRSSList, duplicatedPodcastRSSList := util.RemoveDuplicateItemsInStringSlice(podcastRSSList)

	log.Println(fmt.Sprintf("Loaded %d RSS link(s)", len(podcastRSSList)))

	if len(duplicatedPodcastRSSList) > 0 {
		log.Println(fmt.Sprintf("Found %d duplicate podcast RSS links:", len(duplicatedPodcastRSSList)))
		for index, rss := range duplicatedPodcastRSSList {
			log.Println(fmt.Sprintf("%d. %s", index+1, rss))
		}
	}

	podcastList, failed := podcastParser.ParsePodcastsFromRSSList(podcastRSSList)

	if len(podcastList) != 0 {
		log.Println(fmt.Sprintf("Successfully parsed %d RSS link(s)", len(podcastList)))
	}

	// Print parse failed podcast RSS links
	if len(failed) != 0 {
		log.Println(fmt.Sprintf("%d RSS link(s) parsing failed:", len(failed)))
		for index, rss := range failed {
			log.Println(fmt.Sprintf("%d. %s", index+1, rss))
		}
	}

	// Exit when there are no podcasts to download
	if len(podcastList) == 0 {
		log.Println("No RSS links to download, exit")
		os.Exit(0)
	}

	var podcastDownloadTasks []*podownloader.PodcastDownloadTask
	for _, p := range podcastList {
		tasks := p.GetPodcastDownloadTask(outputFolder, httpClient)
		podcastDownloadTasks = append(podcastDownloadTasks, tasks)
	}
	podcastDownloadTaskIterator := podownloader.NewDownloadTaskIterator(podcastDownloadTasks)
	podcastDownloadTaskIterator.RemoveDownloadedTask(threadCount)

	log.Println("Start download")
	downloadQueue := podownloader.NewDownloadQueueFromDownloadTasks(podcastDownloadTaskIterator.PodcastDownloadTasks)
	failedTaskDestPaths := downloadQueue.StartDownload(httpClient, threadCount)
	log.Println("Download finished")

	// Print failed tasks
	if len(failedTaskDestPaths) > 0 {
		fmt.Println(fmt.Sprintf("%d file(s) download failed:", len(failedTaskDestPaths)))
		for index, failedTaskDestPath := range failedTaskDestPaths {
			fmt.Println(fmt.Sprintf("%d. %s", index+1, failedTaskDestPath))
		}
	}
}

// initConfig initialize configuration items
func initConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else if opmlFilePath == "" && rssListFilePath == "" && rss == "" {
		configPath, err := os.Getwd()
		if err != nil {
			log.Fatalln("Failed to get current directory:", err)
		}
		viper.AddConfigPath(configPath)
		viper.SetConfigName(".podownloader")
	} else {
		// When at least one argument is specified among "opml", "list" and "rss", the configuration file will not be loaded.
		return
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println(fmt.Sprintf("Can't load config file %s: %s", viper.ConfigFileUsed(), err))
		return
	}

	log.Println("Using config file:", viper.ConfigFileUsed())

	rssListFilePath = viper.GetString("list")
	opmlFilePath = viper.GetString("opml")
	rss = viper.GetString("rss")
	outputFolder = viper.GetString("output")
	userAgent = viper.GetString("ua")
	threadCount = viper.GetInt("thread")

	// Print loaded configuration items
	log.Println("Config items:")
	log.Println("RSS list file path:", rssListFilePath)
	log.Println("OPML file path:", opmlFilePath)
	log.Println("RSS:", rss)
	log.Println("Output folder:", outputFolder)
	log.Println("User agent:", userAgent)
	log.Println("Thread count:", threadCount)

	// Exit when no required configuration items in the configuration file
	if opmlFilePath == "" && rssListFilePath == "" && rss == "" {
		log.Fatalln("Please specify at least one argument among \"opml\", \"list\" and \"rss\" in configuration file")
	}
}
