package podownloader

import (
	"PoDownloader/util"
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"net/http"
	"os"
)

// URLDownloadTask is a download task that download a file from URL to Dest
type URLDownloadTask struct {
	JobName string `json:"jobName,omitempty"`
	JobType string `json:"jobType,omitempty"`
	URL     string `json:"url,omitempty"`
	Dest    string `json:"dest,omitempty"`
}

// TextSaveTask is a file save task that save the Text to Dest
type TextSaveTask struct {
	JobName string `json:"jobName,omitempty"`
	JobType string `json:"jobType,omitempty"`
	Text    string `json:"text,omitempty"`
	Dest    string `json:"dest,omitempty"`
}

// EpisodeDownloadTask contains all download tasks in an episode
type EpisodeDownloadTask struct {
	EpisodeTitle           string             `json:"episodeTitle,omitempty"`
	BaseDestDir            string             `json:"baseDestDir,omitempty"`
	EnclosureDownloadTasks []*URLDownloadTask `json:"enclosureDownloadTasks,omitempty"`
	CoverDownloadTask      *URLDownloadTask   `json:"coverDownloadTask,omitempty"`
	ShownotesDownloadTask  *TextSaveTask      `json:"shownotesDownloadTask,omitempty"`
}

// PodcastDownloadTask contains all download tasks in a podcast
type PodcastDownloadTask struct {
	PodcastTitle         string                 `json:"podcastTitle,omitempty"`
	BaseDestDir          string                 `json:"baseDestDir,omitempty"`
	EpisodeDownloadTasks []*EpisodeDownloadTask `json:"episodeDownloadTasks,omitempty"`
	CoverDownloadTask    *URLDownloadTask       `json:"coverDownloadTask,omitempty"`
	RSSDownloadTask      *URLDownloadTask       `json:"rssDownloadTask,omitempty"`
}

// Save writes TextSaveTask.Text to TextSaveTask.Dest
func (t *TextSaveTask) Save() error {
	return util.WriteFile(t.Text, t.Dest)
}

// SaveWithProgress writes TextSaveTask.Text to TextSaveTask.Dest with progress bar
func (t *TextSaveTask) SaveWithProgress(progressBar *mpb.Progress) error {
	taskName := fmt.Sprintf("[Download | %s]", util.FillTextToLength(t.JobType, 9))
	bar := progressBar.AddBar(
		1,
		mpb.PrependDecorators(
			decor.Name(taskName, decor.WC{W: len(taskName) + 1, C: decor.DidentRight}),
			decor.Name(util.GetFirstNCharacters(t.JobName, 20), decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
		),
	)
	err := util.WriteFile(t.Text, t.Dest)
	if err != nil {
		return err
	}
	bar.IncrBy(1)
	return nil
}

// IsDestFileExist returns whether the TextSaveTask destination file exists
func (t *TextSaveTask) IsDestFileExist() bool {
	return util.IsPathExist(t.Dest)
}

// Download downloads URLDownloadTask.URL to URLDownloadTask.Dest
func (c *URLDownloadTask) Download(httpClient *http.Client) error {
	resp, err := httpClient.Get(c.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(c.Dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// DownloadWithProgress downloads URLDownloadTask.URL to URLDownloadTask.Dest with progress bar
func (c *URLDownloadTask) DownloadWithProgress(httpClient *http.Client, progressBar *mpb.Progress) error {
	resp, err := httpClient.Get(c.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(c.Dest)
	if err != nil {
		return err
	}
	defer out.Close()
	taskName := fmt.Sprintf("[Download | %s]", util.FillTextToLength(c.JobType, 9))
	bar := progressBar.AddBar(
		resp.ContentLength,
		mpb.PrependDecorators(
			decor.Name(taskName, decor.WC{W: len(taskName) + 1, C: decor.DidentRight}),
			decor.Name(util.GetFirstNCharacters(c.JobName, 20), decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.EwmaSpeed(decor.UnitKiB, "% .1f", 60),
			decor.Percentage(decor.WC{W: 6}),
		),
	)
	proxyReader := bar.ProxyReader(resp.Body)
	_, err = io.Copy(out, proxyReader)
	return err
}

// IsDestFileExist returns whether the URLDownloadTask destination file exists
func (c *URLDownloadTask) IsDestFileExist() bool {
	return util.IsPathExist(c.Dest)
}

// RemoveDownloadedTask removes all downloaded tasks from EpisodeDownloadTask
func (e *EpisodeDownloadTask) RemoveDownloadedTask() {
	for index, enclosureDownloadTask := range e.EnclosureDownloadTasks {
		if enclosureDownloadTask.IsDestFileExist() {
			e.EnclosureDownloadTasks[index] = nil
		}
	}
	if e.CoverDownloadTask != nil && e.CoverDownloadTask.IsDestFileExist() {
		e.CoverDownloadTask = nil
	}
	if e.ShownotesDownloadTask != nil && e.ShownotesDownloadTask.IsDestFileExist() {
		e.ShownotesDownloadTask = nil
	}
}

// Mkdir creates the episode download destination directory
func (e *EpisodeDownloadTask) Mkdir() error {
	return util.EnsureDirAll(e.BaseDestDir)
}

// RemoveDownloadedTask removes all downloaded cover download tasks from PodcastDownloadTask
func (p *PodcastDownloadTask) RemoveDownloadedTask() {
	if p.CoverDownloadTask != nil && p.CoverDownloadTask.IsDestFileExist() {
		p.CoverDownloadTask = nil
	}
}

// Mkdir creates the podcast download destination directory
func (p *PodcastDownloadTask) Mkdir() error {
	return util.EnsureDirAll(p.BaseDestDir)
}

// RemoveDownloadedTaskAndMakeDirWithProgress removes all downloaded tasks from PodcastDownloadTask
// and creates podcast download destination directory with progress bar
func (p *PodcastDownloadTask) RemoveDownloadedTaskAndMakeDirWithProgress(progressBar *mpb.Progress) {
	taskName := "[Check]"
	job := p.PodcastTitle
	bar := progressBar.AddBar(
		int64(len(p.EpisodeDownloadTasks)),
		mpb.PrependDecorators(
			decor.Name(taskName, decor.WC{W: len(taskName) + 1, C: decor.DidentRight}),
			decor.Name(util.GetFirstNCharacters(job, 15), decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)
	err := p.Mkdir()
	if err != nil {
		return
	}
	p.RemoveDownloadedTask()
	for index, _ := range p.EpisodeDownloadTasks {
		p.EpisodeDownloadTasks[index].Mkdir()
		p.EpisodeDownloadTasks[index].RemoveDownloadedTask()
		bar.IncrBy(1)
	}
}
