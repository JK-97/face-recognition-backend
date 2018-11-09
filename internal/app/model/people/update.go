package people

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

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
