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

// GetPeople gets list of people in db
func GetPeople(ids []string, createdBefore int64, limit int, skip int) ([]*schema.DBPerson, error) {
	opt := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(map[string]int{"created_time": -1})

	filter := make(map[string]interface{})
	if ids != nil {
		filter["_id"] = map[string][]string{"$in": ids}
	}

	if createdBefore != 0 {
		filter["created_time"] = map[string]int64{"$lt": createdBefore}
	}

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
		result = append(result, dbp)
	}

	return result, nil
}
