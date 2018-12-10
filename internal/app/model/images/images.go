package images

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/google/uuid"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func collection() *mongo.Collection {
	return model.DB.Collection("images")
}

// GetImages return a person full images
func GetImages(id string) (*schema.DBImages, error){
	doc := collection().FindOne(context.Background(), map[string]string{"pid": id})
	result := &schema.DBImages{}
	err := doc.Decode(&result)
    if err != nil {
        return nil, err
    }
    return result, nil
}

// UpdateImages update a person's images in db
func UpdateImages(ID string, imgs []string) error {
    updater := map[string]schema.DBImagesUpdater{"$set": schema.DBImagesUpdater{
        PID:     ID,
        Images: imgs,
    }}
    _, err := collection().UpdateOne(context.Background(), map[string]string{"pid": ID}, updater)
    return err
}

// AddImages add a person's image in db
func AddImages(id string, imgs []string) error {
	uid, _ := uuid.NewUUID()
    doc := schema.DBImages{
        ID:             uid.String(),
        PID:            id,
        Images:         imgs,
    }
	_, err := collection().InsertOne(context.Background(), doc)
    return err
}
