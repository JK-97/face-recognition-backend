package people

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// AddPerson add a person to db
func AddPerson(p *schema.Person, images []string) error {
	if len(images) == 0 {
		return fmt.Errorf("should send at least one image")
	}

	uuid, err := uuid.NewUUID()
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
    return err
}

// UpdatePerson update a person in db
func UpdatePerson(p *schema.Person, images []string) error {
	if len(images) == 0 {
		return fmt.Errorf("should send at least one image")
	}

        personID := p.ID
	err := remote.Record(personID, images)
	if err != nil {
		return err
	}

    updater := map[string]*schema.DBPerson{"$set": schema.NewDBPerson(p, images[0])}
	_, err := collection().UpdateOne(context.Background(), map[string]string{"_id": p.ID}, updater)
    return err
}

// DeletePerson delete a person to db
func DeletePerson(id string) error {
	_, err := collection().DeleteOne(context.Background(), map[string]string{"_id": id})
	return err
}
