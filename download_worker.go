package podownloader

import (
	"PoDownloader/logger"
	"context"
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"net/http"
	"sync"
)

type DownloadWorker struct {
	doneWg              *sync.WaitGroup
	httpClient          *http.Client
	progressBar         *mpb.Progress
	logger              *logger.Logger
	IngestChan          chan interface{}
	tasksChan           chan interface{}
	FailedTaskDestPaths []string
	failedTaskListLock  *sync.Mutex
}

func NewDownloadWorker(doneWg *sync.WaitGroup, httpClient *http.Client, progressBar *mpb.Progress, logger *logger.Logger, threadCount int) *DownloadWorker {
	var failedTaskDestPaths []string
	return &DownloadWorker{
		doneWg:              doneWg,
		httpClient:          httpClient,
		progressBar:         progressBar,
		logger:              logger,
		IngestChan:          make(chan interface{}, 1),
		tasksChan:           make(chan interface{}, threadCount),
		FailedTaskDestPaths: failedTaskDestPaths,
		failedTaskListLock:  &sync.Mutex{},
	}
}

func (dw *DownloadWorker) WorkerFunc() {
	defer dw.doneWg.Done()
	for task := range dw.tasksChan {
		if urlDownloadTask, ok := task.(*URLDownloadTask); ok {
			err := urlDownloadTask.DownloadWithProgress(dw.httpClient, dw.progressBar)
			if err != nil {
				dw.logger.Println(fmt.Sprintf("Failed to download %s: %s", urlDownloadTask.URL, err))
				dw.failedTaskListLock.Lock()
				dw.FailedTaskDestPaths = append(dw.FailedTaskDestPaths, urlDownloadTask.Dest)
				dw.failedTaskListLock.Unlock()
			} else {
				dw.logger.PrintlnToFile(fmt.Sprintf("Successfully downloaded %s", urlDownloadTask.Dest))
			}
		} else if textSaveTask, ok := task.(*TextSaveTask); ok {
			err := textSaveTask.SaveWithProgress(dw.progressBar)
			if err != nil {
				dw.logger.Println(fmt.Sprintf("Failed to save %s: %s", textSaveTask.Dest, err))
				dw.failedTaskListLock.Lock()
				dw.FailedTaskDestPaths = append(dw.FailedTaskDestPaths, textSaveTask.Dest)
				dw.failedTaskListLock.Unlock()
			} else {
				dw.logger.PrintlnToFile(fmt.Sprintf("Successfully downloaded %s", textSaveTask.Dest))
			}
		}
	}
}

func (dw *DownloadWorker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for len(dw.tasksChan) > 0 {
				<-dw.tasksChan
			}
			close(dw.tasksChan)
			return
		case job := <-dw.IngestChan:
			dw.tasksChan <- job
		}
	}
}
