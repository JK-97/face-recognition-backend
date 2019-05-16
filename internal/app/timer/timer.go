package timer

import (
    "github.com/jasonlvhit/gocron"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// Callback is timer handler
type Callback func() ()

// UpdateScheduler update task in scheduler
func UpdateScheduler(task Callback, seconds int64) (Callback) {
    gocron.Remove(task)
    gocron.Every(uint64(seconds)).Seconds().Do(task)
    log.Info("UpdateScheduler after ", seconds)
    return task
}

func init() {
    go func() {
        <-gocron.Start()
    } ()
}
