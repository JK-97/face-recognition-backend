package checkin

import (
	"context"

	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

type seal struct {
	startTime int64
	endTime   int64
}

func collection() *mongo.Collection {
	return model.DB.Collection("checkin-history")
}

func saveCheckin(s seal) error {
    l := GetCurrentPeopleSet()
	h := &schema.CheckinHistoryRecord{
        Timestamp:      s.endTime,
        PersonIDS:      l,
	}

	uid, err := uuid.NewUUID()
	if err != nil {
	    _, err = collection().InsertOne(
            context.Background(),
            schema.NewDBCheckinHistoryRecord(h, uid.String()),
        )
	}
	return err
}

// GetHistoryRecords return all person in set during start_time and end_time
func GetHistoryRecords(start_time int64, end_time int64, cameraID string) (*schema.CheckinResp, error) {
	ctx := context.Background()
    cur, err := collection().Find(ctx, bson.D{
        {"timestamp", bson.D{
            {"$gte", start_time},
            {"$lte", end_time},
        }},
    })
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

    ret := &schema.CheckinResp{}
    for pid, _ := range personIDS {
        person, err := people.GetPerson(pid)
        if err != nil {
            return nil, err
        }
        ret.Person = append(ret.Person, person.Person())
    }
    return ret, nil
}
