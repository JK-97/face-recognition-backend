package people

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func collection() *mongo.Collection {
	return model.DB.Collection("people")
}

// CountPeople counts people in db
func CountPeople() (int, error) {
	num, err := collection().Count(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

// GetPerson get a person by id in db
func GetPerson(id string) (*schema.DBPerson, error) {
	doc := collection().FindOne(context.Background(), map[string]string{"_id": id})
	result := &schema.DBPerson{}
	err := doc.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NewFilterPresent creates a filter for present people in checkin
func NewFilterPresent(idsIn []string) map[string]interface{} {
	filter := make(map[string]interface{})
	filter["_id"] = map[string][]string{"$in": idsIn}
	return filter
}

// NewFilterAbsent creates a filter for absent people in checkin
func NewFilterAbsent(idsNotIn []string, createdBefore int64) map[string]interface{} {
	filter := make(map[string]interface{})
	filter["_id"] = map[string][]string{"$nin": idsNotIn}
	filter["created_time"] = map[string]int64{"$lt": createdBefore}
	return filter
}

// GetPeople gets list of people in db
func GetPeople(filter map[string]interface{}, limit int, skip int) ([]*schema.DBPerson, error) {
	opt := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(map[string]int{"created_time": -1})

	ctx := context.Background()
	cur, err := collection().Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	result := []*schema.DBPerson{}
	for cur.Next(ctx) {
		dbp := &schema.DBPerson{}
		if err := cur.Decode(dbp); err != nil {
			return nil, err
		}
		result = append(result, schema.NewDBPersonWithImageURL(dbp))
	}

	return result, nil
}
