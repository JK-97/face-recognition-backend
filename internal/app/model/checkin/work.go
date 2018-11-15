package checkin

import (
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// no sync.locker is needed, because saveCheckin() and checkin() are synced
var currentRecord = schema.CheckinPeopleSet{}

const confThres = 0.5
const countThres = 3

func addRcg(rcg schema.Recognition) {
	log.Debug(rcg)
	if rcg.ID != "unknown" && rcg.Confidence > confThres {
		currentRecord[rcg.ID]++
	}
}

func checkin() error {
	img, err := remote.Capture()
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
}
