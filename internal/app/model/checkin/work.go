package checkin

import (
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

type checkRecordSet map[string]struct{}

// no sync.locker is needed, because saveCheckin() and checkin() are synced
var currentRecord = checkRecordSet{}

func addRcg(rcg schema.Recognition) {
	if model.GetPerson(rcg.ID) == nil {
		log.Debugf("person detected is not in db, id: %v\n", rcg.ID)
		return
	}
	currentRecord[rcg.ID] = struct{}{}
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
