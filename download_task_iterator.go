package podownloader

import (
	"github.com/vbauerster/mpb/v7"
	"sync"
)

// DownloadTaskIterator is the iterator of DownloadTasks
type DownloadTaskIterator struct {
	CurrentIndex         int
	PodcastDownloadTasks []*PodcastDownloadTask
}

// NewDownloadTaskIterator returns *DownloadTaskIterator init from []*PodcastDownloadTask
func NewDownloadTaskIterator(podcastDownloadTasks []*PodcastDownloadTask) *DownloadTaskIterator {
	return &DownloadTaskIterator{
		CurrentIndex:         0,
		PodcastDownloadTasks: podcastDownloadTasks,
	}
}

// Reset resets the current index to 0
func (dti *DownloadTaskIterator) Reset() {
	dti.CurrentIndex = 0
}

// GetLeftLength returns left length of iterator
func (dti *DownloadTaskIterator) GetLeftLength() int {
	return len(dti.PodcastDownloadTasks) - dti.CurrentIndex
}

// Next returns next item of DownloadTaskIterator, and returns nil if next item does not exist
func (dti *DownloadTaskIterator) Next() *PodcastDownloadTask {
	if dti.CurrentIndex < len(dti.PodcastDownloadTasks) {
		podcastDownloadTask := dti.PodcastDownloadTasks[dti.CurrentIndex]
		dti.CurrentIndex += 1
		return podcastDownloadTask
	} else {
		return nil
	}
}

// startRemoveDownloadedTaskAndMakeDir calls PodcastDownloadTask.RemoveDownloadedTaskAndMakeDirWithProgress
// to remove downloaded task and make podcast download destination directory
// and will start a startRemoveDownloadedTaskAndMakeDir goroutine if the iterator has next item
func startRemoveDownloadedTaskAndMakeDir(doneWg *sync.WaitGroup, progressBar *mpb.Progress, task *PodcastDownloadTask, downloadTaskIterator *DownloadTaskIterator) {
	defer doneWg.Done()
	task.RemoveDownloadedTaskAndMakeDirWithProgress(progressBar)
	newTask := downloadTaskIterator.Next()
	if newTask != nil {
		go startRemoveDownloadedTaskAndMakeDir(doneWg, progressBar, newTask, downloadTaskIterator)
	}
}

// RemoveDownloadedTaskAndMakeDir will start ThreadCount startRemoveDownloadedTaskAndMakeDir goroutines to remove
// downloaded tasks and create podcast download destination directory
func (dti *DownloadTaskIterator) RemoveDownloadedTaskAndMakeDir(ThreadCount int) {
	doneWg := new(sync.WaitGroup)
	doneWg.Add(dti.GetLeftLength())
	progressBar := mpb.New(mpb.WithWaitGroup(doneWg))
	for i := 0; i < ThreadCount; i++ {
		task := dti.Next()
		if task != nil {
			go startRemoveDownloadedTaskAndMakeDir(doneWg, progressBar, task, dti)
		}
	}
	progressBar.Wait()
}
