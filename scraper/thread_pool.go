package scraper

import (
	"fmt"
	"time"
)

// job represents the job to be run
type job struct {
	novelName string
}

// worker represents the worker that executes the job
type worker struct {
	scraper Scraper

	workerPool  chan chan job
	jobChannel  chan job
	closeHandle chan bool
}

func newWorker(scraper Scraper, workerPool chan chan job, closeHandle chan bool) *worker {
	fmt.Println("newWorker")
	return &worker{scraper: scraper, workerPool: workerPool, jobChannel: make(chan job), closeHandle: closeHandle}
}

// start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *worker) start() {
	fmt.Println("start")
	go func() {
		for {
			// Put the worker to the worker threadpool
			w.workerPool <- w.jobChannel

			select {
			// Wait for the job
			case j := <-w.jobChannel:
				// Got the job
				w.scraper.ScrapeNovel(j.novelName)
			case <-w.closeHandle:
				// Exit the go routine when the closeHandle channel is closed
				return
			}
		}
	}()
}

// workerPool type for holding the workers and handle the job requests
type workerPool struct {
	scraperName string

	queueSize   int64
	noOfWorkers int

	jobQueue    chan job
	workerPool  chan chan job
	closeHandle chan bool // Channel used to stop all the workers
}

// newWorkerPool creates thread threadpool
func newWorkerPool(scraperName string, noOfWorkers int, queueSize int64) workerPool {
	fmt.Println("newWorkerPool")
	wp := workerPool{scraperName: scraperName, queueSize: queueSize, noOfWorkers: noOfWorkers}
	wp.jobQueue = make(chan job, queueSize)
	wp.workerPool = make(chan chan job, noOfWorkers)
	wp.closeHandle = make(chan bool)
	wp.createPool()
	return wp
}

// createPool creates the workers and start listening on the jobQueue
func (t *workerPool) createPool() error {
	fmt.Println("createPool")
	for i := 0; i < t.noOfWorkers; i++ {
		scraper, err := ScraperCreator(t.scraperName)
		if err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf("failed to create worker pool for %s scraper", t.scraperName)
		}

		worker := newWorker(scraper, t.workerPool, t.closeHandle)
		worker.start()
	}

	fmt.Println("before")
	go t.dispatch()
	time.Sleep(time.Second * 2)
	fmt.Println("after")

	return nil
}

// dispatch listens to the jobqueue and handles the jobs to the workers
func (t *workerPool) dispatch() {
	fmt.Println("dispatch")
	for {
		select {

		case j := <-t.jobQueue:
			fmt.Println("Got job")

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

func (t *workerPool) execute(j job) error {
	// Add the task to the job queue
	if len(t.jobQueue) == int(t.queueSize) {
		return fmt.Errorf("queue is full, not able add the task")
	}
	t.jobQueue <- j
	return nil
}

// Close will close the threadpool
// It sends the stop signal to all the worker that are running
// TODO: need to check the existing /running task before closing the threadpool
func (t *workerPool) close() {
	close(t.closeHandle) // Stops all the routines
	close(t.workerPool)  // Closes the job threadpool
	close(t.jobQueue)    // Closes the job Queue
}

// TODO bad gestion of close need to work on it
