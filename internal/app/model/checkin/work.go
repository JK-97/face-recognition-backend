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

func list(s *schema.CheckinPeopleSet, countThres int) []string {
	l := make([]string, 0)
	for k, v := range *s {
		if v >= countThres {
			l = append(l, k)
		}
	}
	return l
}

func ResetCurrentPeopleSet() {
	recordMutex.Lock()
	currentRecord = schema.CheckinPeopleSet{}
	recordMutex.Unlock()
}

func checkin() {
	devices, _ := device.GetCameras()

	for _, d := range devices {
        img, err := remote.Capture(d.DeviceName)
        if err != nil {
            log.Error("device capture is not working: %c", d.DeviceName)
            continue
        }

        rcgs, err := remote.Detect(img)
        if err != nil {
            log.Error("detect ai is not working: %c", d.DeviceName)
            continue
        }

        for _, rcg := range rcgs {
            addRcg(rcg)
        }
	}
}
