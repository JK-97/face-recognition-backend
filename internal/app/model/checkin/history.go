package checkin

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
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
	expectedCount, err := people.CountPeople()
	if err != nil {
		return err
	}

	h := &schema.CheckinHistory{
		StartTime:     s.startTime,
		EndTime:       s.endTime,
		ExpectedCount: expectedCount,
		ActualCount:   len(currentRecord),
		Record:        currentRecord.List(),
	}
	_, err = collection().InsertOne(context.Background(), h)
	if err != nil {
		return err
	}

	currentRecord = schema.CheckinPeopleSet{}
	return nil
}

// HistoryTimestamps returns all available timestamps for history query
func HistoryTimestamps(limit int, skip int) ([]int64, error) {
	ctx := context.Background()
	opt := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(map[string]int{"start_time": -1}).
		SetProjection(map[string]int{"start_time": 1})
	cur, err := collection().Find(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	result := make([]int64, 0)
	for cur.Next(ctx) {
		t := &struct {
			StartTime int64 `bson:"start_time"`
		}{}
		if err := cur.Decode(t); err != nil {
			return nil, err
		}
		result = append(result, t.StartTime)
	}

	return result, nil
}

// GetHistory returns a checkin record with timestamp
func GetHistory(timestamp int64) (*schema.CheckinHistory, error) {
	doc := collection().FindOne(context.Background(), map[string]int64{"start_time": timestamp})
	result := &schema.CheckinHistory{}
	err := doc.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
