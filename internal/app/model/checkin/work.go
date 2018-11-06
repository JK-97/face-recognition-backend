package checkin

import (
	"fmt"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// no sync.locker is needed, because saveCheckin() and checkin() are synced
var currentRecord = schema.CheckinPeopleSet{}

func addRcg(rcg schema.Recognition) error {
	person, err := people.GetPerson(rcg.ID)
	if person == nil {
		return fmt.Errorf("person detected is not in db, id: %v", rcg.ID)
	}
	if err != nil {
		return err
	}
	currentRecord[rcg.ID] = struct{}{}
	return nil
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
		err := addRcg(rcg)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
