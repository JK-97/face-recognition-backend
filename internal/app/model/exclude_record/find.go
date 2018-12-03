package exclude_record

import (
    "fmt"
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
    "gopkg.in/mgo.v2/bson"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

func collection() *mongo.Collection {
	return model.DB.Collection("exclude_record")
}

// NewFilterExclude creates a filter for present record in exclude_record
// includeBack == true  --> query all time exclude_record
// includeBack == false --> query exclude_record who didn't came back yet
func NewFilterExclude(excludeTime int64, includeBack bool) map[string]interface{} {
	filter := make(map[string]interface{})
	filter["exclude_time"] = map[string]int64{"$lt": excludeTime}
	if !includeBack {
		filter["include_time"] = -1
	}
	return filter
}

// NewFilterExcludeHistory creates a filter for exclude record in exclude_record at timestamp
func NewFilterExcludeHistory(timestamp int64) interface{} {

    const queryStrTem = `{"$and":[{"exclude_time": {"$lt":%d}},{"$or":[{"include_time": {"$eq": -1}},{"include_time": {"$gt": %d}}]}]}`
    queryStr := fmt.Sprintf(queryStrTem, uint64(timestamp), uint64(timestamp))

    // type Et struct {
    //     Gt int64 `json:"$gt" bson:"$gt"`
    // }
    // type It struct {
    //     Lt int64 `json:"$lt" bson:"$lt"`
    //     Eq int64 `json:"$eq" bson:"$eq"`
    // }
    // type HistoryFilter []struct {
    //     ExcludeTime Et     `json:"exclude_time" bson:"exclude_time"`
    //     IncludeTime []It   `json:"include_time" bson:"include_time"`
    // }
    // var vf = HistoryFilter {
    //     {
    //         ExcludeTime: Et {
    //             Gt: timestamp,
    //         },
    //     }, {
    //         IncludeTime:[]It{
    //             It{Lt: timestamp},
    //             It{Eq: -1},
    //         },
    //     },
    // }
    // filter := make(map[string]interface{}) 
    // filter["$and"] = vf

    var filter interface{}
    _ = bson.UnmarshalJSON([]byte(queryStr), &filter)
    return filter
}

// GetExcludeRecord gets list of exclude record in db
func GetExcludeRecord(filter interface{}, limit int, skip int) ([]*schema.DBExcludeRecord, error) {
	ctx := context.Background()
	var cur mongo.Cursor
	var err error
	if limit != -1 && skip != -1 {
		opt := options.Find().
			SetLimit(int64(limit)).
			SetSkip(int64(skip)).
			SetSort(map[string]int{"exclude_time": -1})
		cur, err = collection().Find(ctx, filter, opt)
	} else {
		cur, err = collection().Find(ctx, filter)
	}

	result := []*schema.DBExcludeRecord{}
    if err == mongo.ErrNoDocuments {
        return result, nil
	} else if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		dbp := &schema.DBExcludeRecord{}
		if err := cur.Decode(dbp); err != nil {
			return nil, err
		}
		result = append(result, dbp)
	}

	return result, nil
}

// MakeExcludePeopleSet make set of people
func MakeExcludePeopleSet(records []*schema.DBExcludeRecord) map[string]int64 {
	result := make(map[string]int64)
	for index := 0; index < len(records); index++ {
		excludeTime := records[index].ExcludeTime
		people := records[index].People
		for i := 0; i < len(people); i++ {
			result[people[i].NationalID] = excludeTime
		}
	}
    return result
}

// GetExcludePeopleSet gets set of exclude record in db
func GetExcludePeopleSet(filter interface{}, limit int, skip int) (map[string]int64, error) {
	records, err := GetExcludeRecord(filter, limit, skip)
	if err != nil {
		return nil, err
	}
	return MakeExcludePeopleSet(records), nil
}

// GetExcludePeopleSetNow sample function
func GetExcludePeopleSetNow() (map[string]int64, error) {
	now := util.NowMilli()
	filter := NewFilterExclude(now, false)
	return GetExcludePeopleSet(filter, -1, -1)
}
