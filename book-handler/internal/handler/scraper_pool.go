package handler

import (
	"fmt"
	"time"

	"github.com/MrColorado/backend/book-handler/internal/converter"
	"github.com/MrColorado/backend/book-handler/internal/scraper"
	"github.com/MrColorado/backend/logger"
)

// job represents the job to be run
type job struct {
	NovelName string
}

// worker represents the worker that executes the job
type worker struct {
	scrp  scraper.Scraper
	convs []converter.Converter

	workerPool  chan chan job
	jobChannel  chan job
	closeHandle chan bool
}

func newWorker(scrp scraper.Scraper, convs []converter.Converter, workerPool chan chan job, closeHandle chan bool) *worker {
	logger.Info("newWorker")
	return &worker{
		scrp:        scrp,
		convs:       convs,
		workerPool:  workerPool,
		jobChannel:  make(chan job),
		closeHandle: closeHandle,
	}
}

// start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *worker) start() {
	logger.Info("start")
	go func() {
		for {
			// Put the worker to the worker threadpool
			w.workerPool <- w.jobChannel

			select {
			// Wait for the job
			case j := <-w.jobChannel:
				// Got the job
				w.scrp.ScrapeNovel(j.NovelName)
				for _, conv := range w.convs {
					conv.ConvertNovel(j.NovelName)
				}
			case <-w.closeHandle:
				// Exit the go routine when the closeHandle channel is closed
				return
			}
		}
	}()
}

// workerPool type for holding the workers and handle the job requests
type WorkerPool struct {
	scrpName  string
	convsName []string

	queueSize   int
	noOfWorkers int

	jobQueue    chan job
	workerPool  chan chan job
	closeHandle chan bool // Channel used to stop all the workers
}

// NewWorkerPool creates thread threadpool
func newWorkerPool(scrpName string, convsName []string, noOfWorkers int, queueSize int) WorkerPool {
	logger.Info("NewWorkerPool")

	wp := WorkerPool{
		scrpName:    scrpName,
		convsName:   convsName,
		queueSize:   queueSize,
		noOfWorkers: noOfWorkers,
	}
	wp.jobQueue = make(chan job, queueSize)
	wp.workerPool = make(chan chan job, noOfWorkers)
	wp.closeHandle = make(chan bool)
	wp.createPool()
	return wp
}

// createPool creates the workers and start listening on the jobQueue
func (t *WorkerPool) createPool() error {
	logger.Info("createPool")
	for i := 0; i < t.noOfWorkers; i++ {
		scrp, err := scraper.ScraperCreator(t.scrpName)
		if err != nil {
			logger.Info(err.Error())
			return fmt.Errorf("failed to create worker pool for %s scraper", t.scrpName)
		}

		convs := []converter.Converter{}
		for _, convName := range t.convsName {
			conv, err := converter.ConverterCreator(convName)
			if err != nil {
				logger.Info(err.Error())
				return fmt.Errorf("failed to create conveter %s", convName)
			}
			convs = append(convs, conv)
		}

		worker := newWorker(scrp, convs, t.workerPool, t.closeHandle)
		worker.start()
	}

	logger.Info("before")
	go t.dispatch()
	time.Sleep(time.Second * 2)
	logger.Info("after")

	return nil
}

// dispatch listens to the jobqueue and handles the jobs to the workers
func (t *WorkerPool) dispatch() {
	logger.Info("dispatch")
	for {
		select {

		case j := <-t.jobQueue:
			logger.Info("Got job")

			// Got job
			func(j job) {
				//Find a worker for the job
				jobChannel := <-t.workerPool
				//Submit job to the worker
				jobChannel <- j
			}(j)

		case <-t.closeHandle:
			// Close thread threadpool
			return
		}
	}
}

func (t *WorkerPool) Execute(j job) bool {
	// Add the task to the job queue
	if len(t.jobQueue) == int(t.queueSize) {
		logger.Info("queue is full, not able add the task")
		return false
	}
	t.jobQueue <- j
	return true
}

// Close will close the threadpool
// It sends the stop signal to all the worker that are running
// TODO: need to check the existing /running task before closing the threadpool
func (t *WorkerPool) close() {
	close(t.closeHandle) // Stops all the routines
	close(t.workerPool)  // Closes the job threadpool
	close(t.jobQueue)    // Closes the job Queue
}

// TODO bad gestion of close need to work on it
