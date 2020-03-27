package timer

import (
	"sync"
	"time"
    "sort"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
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

// Push add one job
func (h *JobQueue) Push(x interface{}) {
	timerMutex.Lock()
	*h = append(*h, x.(*Job))
	timerMutex.Unlock()
}

// Pop delete one job
func (h *JobQueue) Pop() interface{} {
	timerMutex.Lock()
	old := *h
	x := old[0]
	*h = old[1:]
	timerMutex.Unlock()
	return x
}

// NextTimestamp get last timestamp in job queue
func (h *JobQueue) NextTimestamp() time.Duration {
	return time.Duration((*h)[0].timestamp)
}

var checkTimer *time.Timer
var isRunTimer = false
var handlerQueue = &JobQueue{}
var timerMutex = &sync.Mutex{}

// Init Timer
func Init() {
    InitReload()
    for _, i := range *handlerQueue {
        i.timestamp, _ = i.cb(true)
    }
    RunNextTimer()
}


// RunNextTimer run timer
func RunNextTimer() {

    // TODO 二分查找
	timerMutex.Lock()
    sort.Sort(handlerQueue)
	timerMutex.Unlock()

	if len(*handlerQueue) != 0 {

        last := int64(handlerQueue.NextTimestamp())
        waiting := time.Duration(last - util.NowMilli())

        if waiting <= 0 {
			DoJob()
        } else {
            checkTimer = time.NewTimer(time.Millisecond * waiting)
            go func() {
                <-checkTimer.C
                timerMutex.Lock()
                isRunTimer = true
                timerMutex.Unlock()
                DoJob()
                timerMutex.Lock()
                isRunTimer = false
                timerMutex.Unlock()
            }()
        }

	} else {
		log.Info("Timer won't run: no job")
	}
}

// DoJob do timer handle
func DoJob() {
	job := handlerQueue.Pop().(*Job)

	nextTime, err := job.cb(false)
	if err != nil {
		log.Error(err)
	}

	if nextTime != 0 {
        job.timestamp = nextTime
		handlerQueue.Push(job)
    }

	RunNextTimer()
}

// RegisterHandler registe timer handler
func RegisterHandler(cb Callback) *Job {
    var job = &Job{
		cb:        cb,
		timestamp: -1,
    }

	handlerQueue.Push(job)
	return job
}

// UpdateTimer stop timer and retry
func UpdateTimer(job *Job, nextTime int64) {
	log.Info("Timer will stop: try to update")
    timerMutex.Lock()
    if isRunTimer {
	    checkTimer.Stop()
    }
    timerMutex.Unlock()

    if nextTime != 0 {
        job.timestamp = nextTime
	    handlerQueue.Push(job)
    }

	RunNextTimer()
}
