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

func GetCurrentPeopleSet() []string {
	recordMutex.Lock()
	l := list(&currentRecord, countThres)
	recordMutex.Unlock()
	return l
}

// LoadHistoryResult load history checkin recordset into memory
func LoadHistoryResult(t int64) error {
	h, err := GetHistory(t)
	if err != nil {
		return err
	}
	recordMutex.Lock()
	for _, v := range h.Record {
		currentRecord[v] = countThres
	}
	recordMutex.Unlock()
	return nil
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
