package checkin

import (
	"sync"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/device"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// need lock by mutli-camera access
var currentRecord = schema.CheckinPeopleSet{}
var recordMutex = &sync.Mutex{}

const confThres = 0.5
const countThres = 3

func addRcg(rcg schema.Recognition) {
	log.Debug(rcg)
	if rcg.ID != "unknown" && rcg.Confidence > confThres {
		recordMutex.Lock()
		currentRecord[rcg.ID]++
		recordMutex.Unlock()
	}
}

func checkin() error {

	devices, _ := device.GetCameras()

	for _, d := range devices {
		go func(d *schema.Camera) error {
			img, err := remote.Capture(d.DeviceName)
			if err != nil {
				return err
			}

			rcgs, err := remote.Detect(img)
			if err != nil {
				return err
			}

			for _, rcg := range rcgs {
				addRcg(rcg)
			}
			return nil
		}(d)
	}

	return nil
}
