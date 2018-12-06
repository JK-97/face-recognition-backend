package timer

import (
	"container/heap"
	"sync"
	"time"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// Callback is timer handler
type Callback func(bool) (int64, error)

// Job is warpper of timer handler
type Job struct {
	cb        Callback
	timestamp int64
}

// JobQueue is job's min-heap
type JobQueue []*Job

func (h JobQueue) Len() int           { return len(h) }
func (h JobQueue) Less(i, j int) bool { return h[i].timestamp < h[j].timestamp }
func (h JobQueue) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push add one job into heap
func (h *JobQueue) Push(x interface{}) {
	*h = append(*h, x.(*Job))
}

// Pop delete one job into heap
func (h *JobQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// LastTimestamp get last timestamp in job queue
func (h *JobQueue) LastTimestamp() time.Duration {
	return time.Duration((*h)[len(*h)-1].timestamp)
}

var checkTimer *time.Timer
var handlerQueue = &JobQueue{}
var timerMutex = &sync.Mutex{}

// RunNextTimer run timer
func RunNextTimer() {
	if len(*handlerQueue) != 0 {

        for handlerQueue.LastTimestamp() == -1 {
            timerMutex.Lock()
            job := heap.Pop(handlerQueue).(*Job)
            timerMutex.Unlock()
            job.timestamp, _ = job.cb(true)
            timerMutex.Lock()
		    handlerQueue.Push(job)
            timerMutex.Unlock()
        }

		checkTimer = time.NewTimer(time.Second * handlerQueue.LastTimestamp())
		go func() {
			<-checkTimer.C
			DoJob()
		}()

	} else {
		log.Info("Timer won't run: no job")
	}
}

// DoJob do timer handle
func DoJob() {
	timerMutex.Lock()
	job := heap.Pop(handlerQueue).(*Job)
	timerMutex.Unlock()

	nextTime, err := job.cb(false)
	if err != nil {
		log.Error(err)
	}

	if nextTime != 0 {
        job.timestamp = nextTime
	    timerMutex.Lock()
		handlerQueue.Push(job)
	    timerMutex.Unlock()
    }

	RunNextTimer()
}

// RegisterHandler registe timer handler
func RegisterHandler(cb Callback) *Job {
	timerMutex.Lock()
	if handlerQueue == nil {
		heap.Init(handlerQueue)
	}
	timerMutex.Unlock()

    var job = &Job{
		cb:        cb,
		timestamp: -1,
    }

	timerMutex.Lock()
	handlerQueue.Push(job)
	timerMutex.Unlock()

	return job
}

// UpdateTimer stop timer and retry
func UpdateTimer(job *Job, nextTime int64) {
	log.Info("Timer will stop: try to update")
	checkTimer.Stop()

    // TODO 二分查找
	timerMutex.Lock()
    if nextTime != 0 {
        job.timestamp = nextTime
	    handlerQueue.Push(job)
    }
	heap.Init(handlerQueue)
	timerMutex.Unlock()

	RunNextTimer()
}
