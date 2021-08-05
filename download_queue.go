package podownloader

import (
	"PoDownloader/logger"
	"context"
	"errors"
	"github.com/vbauerster/mpb/v7"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// DownloadQueue is the queue of download tasks
// include two types of download tasks: URLDownloadTask and TextSaveTask
type DownloadQueue struct {
	tasks []interface{}
	lock  *sync.Mutex
}

// NewDownloadQueueFromDownloadTasks converts []*PodcastDownloadTask to *DownloadQueue
// and returns the converted *DownloadQueue
// *DownloadQueue will contains 5 types of download tasks:
// 1. Podcast cover download task
// 2. Podcast RSS download task
// 3. Episode cover download task
// 4. Episode shownotes download task
// 5. Episodes enclosures download task
// All nil tasks will be filtered out
func NewDownloadQueueFromDownloadTasks(podcastDownloadTasks []*PodcastDownloadTask) *DownloadQueue {
	var tasks []interface{}
	for _, podcastDownloadTask := range podcastDownloadTasks {
		tasks = append(tasks, podcastDownloadTask.RSSDownloadTask)
		if podcastDownloadTask.CoverDownloadTask != nil {
			tasks = append(tasks, podcastDownloadTask.CoverDownloadTask)
		}
		for _, episodeDownloadTask := range podcastDownloadTask.EpisodeDownloadTasks {
			if episodeDownloadTask.ShownotesDownloadTask != nil {
				tasks = append(tasks, episodeDownloadTask.ShownotesDownloadTask)
			}
			if episodeDownloadTask.CoverDownloadTask != nil {
				tasks = append(tasks, episodeDownloadTask.CoverDownloadTask)
			}
			for _, enclosureDownloadTask := range episodeDownloadTask.EnclosureDownloadTasks {
				if enclosureDownloadTask != nil {
					tasks = append(tasks, enclosureDownloadTask)
				}
			}
		}
	}
	return &DownloadQueue{
		tasks: tasks,
		lock:  &sync.Mutex{},
	}
}

// EnQueue adds an element to the rear of the queue
func (dq *DownloadQueue) EnQueue(podcastDownloadTasks *PodcastDownloadTask) {
	dq.lock.Lock()
	dq.tasks = append(dq.tasks, podcastDownloadTasks)
	dq.lock.Unlock()
}

// DeQueue removes an element from the front of the queue
func (dq *DownloadQueue) DeQueue() (interface{}, error) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	if len(dq.tasks) > 0 {
		frontDownloadTask := dq.tasks[0]
		dq.tasks = dq.tasks[1:]
		return frontDownloadTask, nil
	}
	return nil, errors.New("queue is empty")
}

// Front returns queue front
func (dq *DownloadQueue) Front() (interface{}, error) {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	if len(dq.tasks) > 0 {
		return dq.tasks[0], nil
	}
	return nil, errors.New("queue is empty")
}

// Length returns queue length
func (dq *DownloadQueue) Length() int {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	return len(dq.tasks)
}

// IsEmpty returns whether the queue is empty
func (dq *DownloadQueue) IsEmpty() bool {
	dq.lock.Lock()
	defer dq.lock.Unlock()
	return len(dq.tasks) == 0
}

func (dq *DownloadQueue) StartDownload(threadCount int, httpClient *http.Client, logger *logger.Logger) []string {
	taskCount := dq.Length()
	doneWg := new(sync.WaitGroup)
	doneWg.Add(threadCount)
	progressBar := mpb.New(
		mpb.WithWaitGroup(doneWg),
	)
	ctx, cancelFunc := context.WithCancel(context.Background())
	downloadWorker := NewDownloadWorker(doneWg, httpClient, progressBar, logger, threadCount)

	go downloadWorker.Start(ctx)
	for i := 0; i < threadCount; i++ {
		go downloadWorker.WorkerFunc()
	}

	go func() {
		termChan := make(chan os.Signal)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

		<-termChan

		logger.Println("Received cancellation signal, waiting for all download tasks done")
		cancelFunc()
	}()

	go func() {
		for i := 0; i < taskCount; i++ {
			task, err := dq.DeQueue()
			if err != nil {

			}
			downloadWorker.IngestChan <- task
		}
	}()

	doneWg.Wait()
	return downloadWorker.FailedTaskDestPaths
}
