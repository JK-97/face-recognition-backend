package checkin

import (
	"sync"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/device"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// need lock by mutli-camera access
var currentRecord = schema.DeviceCheckinPeopleSet{}
var recordMutex = &sync.Mutex{}

const confThres = 0.5
const countThres = 3

func addRcg(rcg schema.Recognition, cameraID string) {
	log.Debug(rcg)
	if rcg.ID != "unknown" && rcg.Confidence > confThres {
		recordMutex.Lock()
        if _, ok := currentRecord[cameraID]; !ok {
            currentRecord[cameraID] = map[string]int{}
        }
		currentRecord[cameraID][rcg.ID]++
		recordMutex.Unlock()
	}
}

func list(s *schema.DeviceCheckinPeopleSet, countThres int, cameraID string) []string {
	l := make([]string, 0)

    var recordSet = &schema.DeviceCheckinPeopleSet{}
    if cameraID == "" {
        recordSet = s
    } else if v, ok := (*s)[cameraID]; ok {
        (*recordSet)[cameraID] = v
    }

	for _, people_set := range *s {
        for k, v := range people_set {
            if v >= countThres {
                l = append(l, k)
            }
        }
	}
	return l
}

func listDevices(s *schema.DeviceCheckinPeopleSet) []string {
	l := make([]string, 0)
	for device := range *s {
        l = append(l, device)
    }
    return l
}

func ResetCurrentPeopleSet() {
	recordMutex.Lock()
	currentRecord = schema.DeviceCheckinPeopleSet{}
	recordMutex.Unlock()
}

func GetCurrentDevices() []string {
	recordMutex.Lock()
	l := listDevices(&currentRecord)
	recordMutex.Unlock()
    return l
}

func GetCurrentPeopleSet(cameraID string) []string {
	recordMutex.Lock()
	l := list(&currentRecord, countThres, cameraID)
	recordMutex.Unlock()
	return l
}

func checkin() {
	devices, _ := device.GetCameras()

	for _, d := range devices {
        img, err := remote.Capture(d.CameraID)
        if err != nil {
            log.Error("device capture is not working: %c", d.CameraID)
            continue
        }

        rcgs, err := remote.Detect(img)
        if err != nil {
            log.Error("detect ai is not working: %c", d.CameraID)
            continue
        }

        for _, rcg := range rcgs {
            addRcg(rcg, d.CameraID)
        }
	}
}
