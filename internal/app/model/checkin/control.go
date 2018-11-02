package checkin

import (
	"fmt"
	"sync"
	"time"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

const (
	checking schema.StatusOption = "checking"
	stopped  schema.StatusOption = "stopped"
)

var status = stopped
var checkinSeal = make(chan seal)
var statusMutex = sync.Mutex{}
var startTime int64

// Status returns checkin status
func Status() schema.StatusOption {
	return status
}

// changeStatus changes status, return true if the status is changed
func changeStatus(newStatus schema.StatusOption) bool {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	if status == newStatus {
		return false
	}

	status = newStatus
	return true
}

// StartCheckin starts check in, return true if the status is changed
func StartCheckin() error {
	changed := changeStatus(checking)
	if changed == false {
		return fmt.Errorf("checkin already started")
	}
	startTime = util.NowMilli()
	go func() {
		for {
			select {
			case s := <-checkinSeal:
				saveCheckin(s)
				return
			default:
				err := checkin()
				if err != nil {
					log.Error(err)
				}
				time.Sleep(time.Millisecond * 300)
			}
		}
	}()
	return nil
}

// StopCheckin stops check in, return true if the status is changed
func StopCheckin() (int64, error) {
	changed := changeStatus(stopped)
	if changed == false {
		return 0, fmt.Errorf("checkin is not started")
	}
	t := util.NowMilli()
	checkinSeal <- seal{
		startTime: startTime,
		endTime:   t,
	}
	return startTime, nil
}
