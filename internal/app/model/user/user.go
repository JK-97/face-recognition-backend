package user

import (
	"context"
	"errors"

	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// ErrUserNotFound not found user in db
var ErrUserNotFound = errors.New("User not found")

func collection() *mongo.Collection {
	return model.DB.Collection("user")
}

// FindUser get cameralist in db
func FindUser(userName string) (*schema.User, error) {
	doc := collection().FindOne(context.Background(), map[string]string{"username": userName})
	result := &schema.User{}
	err := doc.Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return result, nil
}
