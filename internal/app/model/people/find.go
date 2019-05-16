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

// GetPerson get a person by id in db
func GetPerson(filter map[string]string) (*schema.DBPerson, error) {
	doc := collection().FindOne(context.Background(), filter)
	result := &schema.DBPerson{}
	err := doc.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PersonFilter return filter of db person
func PersonFilter(id string, serialNumber string) map[string]string {
    var ret map[string]string
    if id != "" {
        ret["_id"] = id
    }
    if serialNumber != "" {
        ret["serial_number"] = serialNumber
    }
    return ret
}

func GetPeople(filter map[string]interface{}, limit int, skip int) ([]*schema.DBPerson, error) {
	opt := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(map[string]int{"serial_number": -1})

	ctx := context.Background()
	cur, err := collection().Find(ctx, filter, opt)

	result := []*schema.DBPerson{}
    if err == mongo.ErrNoDocuments {
        return result, nil
	} else if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		dbp := &schema.DBPerson{}
		if err := cur.Decode(dbp); err != nil {
			return nil, err
		}
		result = append(result, dbp)
	}

	return result, nil
}
