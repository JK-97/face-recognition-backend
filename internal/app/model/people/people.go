package people

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"github.com/satori/go.uuid"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
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
func GetPeople(limit int, skip int) ([]*schema.DBPerson, error) {
	ctx := context.Background()
	opt := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(map[string]int{"created_time": -1})
	cur, err := collection().Find(ctx, nil, opt)
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

// AddPerson add a person to db
func AddPerson(p *schema.Person, images []string) error {
	if len(images) == 0 {
		return fmt.Errorf("should send at least one image")
	}

	uuid, err := uuid.NewV1()
	if err != nil {
		return err
	}
	personID := uuid.String()

	err = remote.Record(personID, images)
	if err != nil {
		return err
	}

	p.ID = personID
	dbp := schema.NewDBPerson(p, images[0])

	_, err = collection().InsertOne(context.Background(), dbp)
	if err != nil {
		return err
	}

	return nil
}

// DeletePerson delete a person to db
func DeletePerson(id string) error {
	_, err := collection().DeleteOne(context.Background(), map[string]string{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
