package podcast

import (
	podownloader "PoDownloader"
	"PoDownloader/logger"
	"PoDownloader/util"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

// Podcast contains all information about a podcast
type Podcast struct {
	RSS         string               `json:"rss,omitempty"`
	Title       string               `json:"title,omitempty"`
	SafeTitle   string               `json:"safeTitle,omitempty"`
	Description string               `json:"description,omitempty"`
	ITunesExt   *ITunesFeedExtension `json:"iTunesExt,omitempty"`
	Items       []*Item              `json:"items,omitempty"`
}

// ITunesFeedExtension is the extension fields of Podcast
type ITunesFeedExtension struct {
	Author     string       `json:"author,omitempty"`
	Categories []*Category  `json:"categories,omitempty"`
	Owner      *ITunesOwner `json:"owner,omitempty"`
	Subtitle   string       `json:"subtitle,omitempty"`
	Summary    string       `json:"summary,omitempty"`
	Image      string       `json:"image,omitempty"`
	Explicit   string       `json:"explicit,omitempty"`
}

// ITunesOwner is the Owner field of ITunesFeedExtension
type ITunesOwner struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// GetItemCount returns the number of items in a Podcast instance
func (p *Podcast) GetItemCount() int {
	return len(p.Items)
}

// GetPodcastDownloadDestDir returns the podcast download destination directory path
// Download dir = output dir + podcast title
func (p *Podcast) GetPodcastDownloadDestDir(destDir string) string {
	return path.Join(destDir, p.SafeTitle)
}

// GetPodcastDownloadTask returns a podownloader.PodcastDownloadTask instance from a Podcast instance
func (p *Podcast) GetPodcastDownloadTask(destDir string, httpClient *http.Client, logger *logger.Logger) *podownloader.PodcastDownloadTask {
	podcastDownloadDestDir := p.GetPodcastDownloadDestDir(destDir)

	// Podcast cover download task
	var podcastCoverDownloadTask *podownloader.URLDownloadTask = nil
	if p.ITunesExt.Image != "" {
		podcastCoverExtensionName, err := util.GetRemoteFileExtensionName(httpClient, p.ITunesExt.Image)
		if err != nil {
			logger.Println(fmt.Sprintf("Failed to get cover extension name of podcast [%s]: %s", p.Title, p.ITunesExt.Image))
		}
		podcastCoverDownloadDest := path.Join(podcastDownloadDestDir, fmt.Sprintf("cover.%s", podcastCoverExtensionName))
		podcastCoverDownloadTask = &podownloader.URLDownloadTask{
			JobName: p.Title,
			JobType: "Cover",
			URL:     p.ITunesExt.Image,
			Dest:    podcastCoverDownloadDest,
		}
	}

	// Episode download task
	var episodeDownloadTasks []*podownloader.EpisodeDownloadTask
	for _, item := range p.Items {
		// item dest dir = download dir + episode title
		itemDownloadDestDir := item.GetItemDownloadDestDir(podcastDownloadDestDir)

		// Cover download task
		var episodeCoverDownloadTask *podownloader.URLDownloadTask = nil
		if item.ITunesExt.Image != "" {
			episodeCoverExtensionName, err := util.GetRemoteFileExtensionName(httpClient, item.ITunesExt.Image)
			if err != nil {
				logger.Println(fmt.Sprintf("Failed to get cover extension name of episode [%s] - [%s]: %s", p.Title, item.Title, item.ITunesExt.Image))
			} else {
				episodeCoverDownloadTask = &podownloader.URLDownloadTask{
					JobName: fmt.Sprintf("%s - %s", p.Title, item.Title),
					JobType: "Cover",
					URL:     item.ITunesExt.Image,
					Dest:    path.Join(itemDownloadDestDir, fmt.Sprintf("cover.%s", episodeCoverExtensionName)),
				}
			}
		}

		// Shownotes download task
		var shownoteDownloadTask *podownloader.TextSaveTask
		if item.Description != "" {
			shownoteDownloadTask = &podownloader.TextSaveTask{
				JobName: fmt.Sprintf("%s - %s", p.Title, item.Title),
				JobType: "Shownotes",
				Text:    item.Description,
				Dest:    path.Join(itemDownloadDestDir, "shownotes.html"),
			}
		}

		// Enclosure download task
		var enclosureDownloadTasks []*podownloader.URLDownloadTask
		for index, enclosure := range item.Enclosures {
			enclosureExtensionName, err := enclosure.GetEnclosureFileExtensionName(httpClient)
			if err != nil {
				logger.Println(fmt.Sprintf("Failed to get enclosure extension name of [%s] - [%s]: %s", p.Title, item.Title, enclosure.URL))
			} else {
				var (
					enclosureFileName string
					jobName           string
				)
				if len(item.Enclosures) == 1 {
					// Only one enclosure, no need to append enclosure index to the file name
					enclosureFileName = fmt.Sprintf("%s.%s", item.SafeTitle, enclosureExtensionName)
					jobName = fmt.Sprintf("%s - %s", p.Title, item.Title)
				} else {
					// More than one enclosure, need to append enclosure index to the file name
					enclosureFileName = fmt.Sprintf("%s_%d.%s", item.SafeTitle, index+1, enclosureExtensionName)
					jobName = fmt.Sprintf("%s - %s #%d", p.Title, item.Title, index+1)
				}
				enclosureDownloadTasks = append(enclosureDownloadTasks, &podownloader.URLDownloadTask{
					JobName: jobName,
					JobType: "Enclosure",
					URL:     enclosure.URL,
					Dest:    path.Join(itemDownloadDestDir, enclosureFileName),
				})
			}
		}

		episodeDownloadTasks = append(episodeDownloadTasks, &podownloader.EpisodeDownloadTask{
			EpisodeTitle:           item.Title,
			BaseDestDir:            itemDownloadDestDir,
			EnclosureDownloadTasks: enclosureDownloadTasks,
			CoverDownloadTask:      episodeCoverDownloadTask,
			ShownotesDownloadTask:  shownoteDownloadTask,
		})
	}

	return &podownloader.PodcastDownloadTask{
		PodcastTitle:         p.Title,
		BaseDestDir:          podcastDownloadDestDir,
		EpisodeDownloadTasks: episodeDownloadTasks,
		CoverDownloadTask:    podcastCoverDownloadTask,
		RSSDownloadTask: &podownloader.URLDownloadTask{
			JobName: fmt.Sprintf("%s | RSS", p.Title),
			JobType: "RSS",
			URL:     p.RSS,
			Dest:    path.Join(podcastDownloadDestDir, "rss.xml"),
		},
	}
}

// GetJSON returns a Podcast instance JSON format
func (p *Podcast) GetJSON() (string, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
