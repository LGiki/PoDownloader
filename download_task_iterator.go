package podownloader

import (
	"github.com/vbauerster/mpb/v8"
	"sync"
)

// DownloadTaskIterator is the iterator of DownloadTasks
type DownloadTaskIterator struct {
	CurrentIndex         int
	PodcastDownloadTasks []*PodcastDownloadTask
	lock                 *sync.Mutex
}

// NewDownloadTaskIterator returns *DownloadTaskIterator init from []*PodcastDownloadTask
func NewDownloadTaskIterator(podcastDownloadTasks []*PodcastDownloadTask) *DownloadTaskIterator {
	return &DownloadTaskIterator{
		CurrentIndex:         0,
		PodcastDownloadTasks: podcastDownloadTasks,
		lock:                 &sync.Mutex{},
	}
}

// Reset resets the current index to 0
func (dti *DownloadTaskIterator) Reset() {
	dti.lock.Lock()
	dti.CurrentIndex = 0
	dti.lock.Unlock()
}

// GetLeftLength returns left length of iterator
func (dti *DownloadTaskIterator) GetLeftLength() int {
	dti.lock.Lock()
	defer dti.lock.Unlock()
	return len(dti.PodcastDownloadTasks) - dti.CurrentIndex
}

// Next returns the next item of DownloadTaskIterator, and returns nil if next item does not exist
func (dti *DownloadTaskIterator) Next() *PodcastDownloadTask {
	dti.lock.Lock()
	defer dti.lock.Unlock()
	if dti.CurrentIndex < len(dti.PodcastDownloadTasks) {
		podcastDownloadTask := dti.PodcastDownloadTasks[dti.CurrentIndex]
		dti.CurrentIndex++
		return podcastDownloadTask
	}
	return nil
}

// startRemoveDownloadedTask calls PodcastDownloadTask.RemoveDownloadedTaskWithProgress
// to remove downloaded task and will start a startRemoveDownloadedTask goroutine
// if the iterator has next item
func startRemoveDownloadedTask(doneWg *sync.WaitGroup, progressBar *mpb.Progress, task *PodcastDownloadTask, downloadTaskIterator *DownloadTaskIterator) {
	defer doneWg.Done()
	task.RemoveDownloadedTaskWithProgress(progressBar)
	newTask := downloadTaskIterator.Next()
	if newTask != nil {
		go startRemoveDownloadedTask(doneWg, progressBar, newTask, downloadTaskIterator)
	}
}

// RemoveDownloadedTask will start ThreadCount startRemoveDownloadedTask goroutines to remove
// downloaded tasks
func (dti *DownloadTaskIterator) RemoveDownloadedTask(threadCount int) {
	doneWg := new(sync.WaitGroup)
	doneWg.Add(dti.GetLeftLength())
	progressBar := mpb.New(mpb.WithWaitGroup(doneWg))
	for i := 0; i < threadCount; i++ {
		task := dti.Next()
		if task != nil {
			go startRemoveDownloadedTask(doneWg, progressBar, task, dti)
		}
	}
	progressBar.Wait()
}
