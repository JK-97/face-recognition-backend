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

var userInited = false
func initUser(){
	if !userInited {
		userInited = true
		num, err := collection().Count(context.Background(), map[string]string{})
		if err == nil && num == 0 {
			collection().InsertOne(context.Background(), map[string]string{"username": "admin", "password": "e10adc3949ba59abbe56e057f20f883e", "_id": "admin"})
		}
	}
}

// FindUser get users in db
func FindUser(userName string) (*schema.User, error) {
	initUser()
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
