package images

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func collection() *mongo.Collection {
	return model.DB.Collection("images")
}

// GetImages return a person full images
func GetImages(id string) (*schema.DBImages, error){
	doc := collection().FindOne(context.Background(), map[string]string{"national_id": id})
	result := &schema.DBImages{}
	err := doc.Decode(&result)
    if err != nil {
        return nil, err
    }
    return result, nil
}

// UpdateImages update a person's images in db
func UpdateImages(nationalID string, imgs []string) error {
    updater := schema.DBImages{
        NationalID:     nationalID,
        Images:         imgs,
    }
	_, err := collection().UpdateOne(context.Background(), map[string]string{"national_id": nationalID}, updater)
    return err
}

// AddImages add a person's image in db
func AddImages(nationalID string, imgs []string) error {
    doc := schema.DBImages{
        NationalID:     nationalID,
        Images:         imgs,
    }
	_, err := collection().InsertOne(context.Background(), doc)
    return err
}
