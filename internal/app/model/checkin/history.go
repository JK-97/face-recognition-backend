package checkin

import (
	"context"

	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/mongo"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

type seal struct {
	startTime int64
	endTime   int64
}

func collection() *mongo.Collection {
	return model.DB.Collection("checkin-history")
}

func saveCheckinWithDevice(cameraID string, endTime int64) error {
    l := GetCurrentPeopleSet(cameraID)
	h := &schema.CheckinHistoryRecord{
        Timestamp:      endTime,
        PersonIDS:      l,
        CameraID:       cameraID,
	}

    log.Info("saveCheckinWithDevice ", cameraID, ":", endTime, ":", l)
	uid, err := uuid.NewUUID()
	if err == nil {
	    _, err = collection().InsertOne(
            context.Background(),
            schema.NewDBCheckinHistoryRecord(h, uid.String()),
        )
	}
	return err
}

func saveCheckin(s seal) error {
    devices := GetCurrentDevices()
    for _, device := range devices {
        err := saveCheckinWithDevice(device, s.endTime)
        if err != nil {
            return err
        }
    }
    return nil
}

// GetHistoryRecords return all person in set during start_time and end_time
func GetHistoryRecords(start_time int64, end_time int64, cameraID string) (*schema.CheckinResp, error) {
	ctx := context.Background()


    cond := make(map[string]interface{})
    cond["timestamp"] = map[string]int64 {
        "$gte": start_time,
        "$lte": end_time,
    }
    if cameraID != "" {
        cond["camera_id"] = cameraID
    }

    cur, err := collection().Find(ctx, cond)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

    personIDS := schema.CheckinPeopleSet{}
	for cur.Next(ctx) {
        r := &schema.DBCheckinHistoryRecord{}
		if err := cur.Decode(r); err != nil {
			return nil, err
		}

        for _, person := range r.PersonIDS {
            personIDS[person] = 1
        }
	}

    ret := &schema.CheckinResp{
        StartTime: start_time,
        EndTime: end_time,
        Person: make([]*schema.Person, 0),
    }
    for pid, _ := range personIDS {
        person, err := people.GetPerson(people.PersonFilter(pid, ""))
        if err != nil {
            return nil, err
        }
        ret.Person = append(ret.Person, person.Person())
    }
    return ret, nil
}
